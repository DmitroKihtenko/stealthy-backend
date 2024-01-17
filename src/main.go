package main

import (
	"SharingBackend/api"
	"SharingBackend/api/controllers"
	"SharingBackend/api/services"
	"SharingBackend/base"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
	"time"

	"SharingBackend/docs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title Sharing Backend
// @version 1.0.0
// @schemes http
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey User
// @In header
// @Name Access token
// @Description Access token of authenticated user

func createMongoClient(config *base.SharingBackendConfig, ctx *context.Context) *mongo.Client {
	timeout := time.Duration(config.MongoDB.SecondsTimeout) * time.Second
	clientOptions := options.Client().
		ApplyURI(config.MongoDB.URL).
		SetConnectTimeout(timeout).
		SetSocketTimeout(timeout).
		SetServerSelectionTimeout(timeout).
		SetTimeout(timeout)

	client, err := mongo.Connect(*ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	return client
}

func checkMongoConnection(client *mongo.Client, ctx *context.Context) {
	base.Logger.Info("Checking mongo DB connection")
	if err := client.Ping(*ctx, nil); err != nil {
		panic(err)
	}
}

func processError(err error) {
	base.Logger.WithFields(logrus.Fields{
		"error": err.Error(),
	}).Fatal("Exiting due to fatal error")
	os.Exit(1)
}

func processPanic() {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if ok {
			base.Logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Fatal("Exiting due to fatal error")
		} else {
			base.Logger.WithFields(logrus.Fields{
				"detail": r,
			}).Fatal("Exiting due to fatal error")
		}
		os.Exit(1)
	}
}

func configureSwagger(swaggerRouter *gin.RouterGroup, config *base.SharingBackendConfig) {
	base.Logger.Info("Configuring openapi")

	baseUrl := config.Server.Socket + config.Server.BasePath

	docs.SwaggerInfo.Host = config.Server.Socket
	docs.SwaggerInfo.BasePath = config.Server.BasePath
	docs.SwaggerInfo.Description = "Stealthy backend service. " +
		"REST API web application. Encapsulates user's service " +
		"business logic of Stealthy system." +
		"<br><br>API is based on JSON (JavaScript Object Notation) Web " +
		"Application Services and HTTPS transport, so is accessible from " +
		"any platform or operating system. Connection to the JSON API is " +
		"provided via HTTPS. Authorization is performed using an " +
		"authorization access token. The example below illustrates " +
		"\"Get authorized user data\" request with an access token: " +
		"<br><strong>curl -X GET " + baseUrl + "/v1/users/me " +
		"-H \"Authorization: Bearer " +
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.ey\"</strong>" +
		"<h4>How To Get Authorization Data</h4>" +
		"An authorization access token can be received using request \"" +
		"Sign-in user\". The example below illustrates receiving this token " +
		"with the cURL CLI tool: <br><strong>curl -X POST --data-binary " +
		"'{\"password\": \"p@ssw0rd\",\"username\": \"john_doe\"}' " +
		baseUrl + "/v1/login</strong>"

	swaggerRouter.GET(
		config.Server.OpenapiBasePath+"/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)
}

func setLogger(config *base.SharingBackendConfig) {
	base.Logger = base.CreateLogger(config)
}

func runServer(engine *gin.Engine, config *base.SharingBackendConfig) {
	base.Logger.Info("Starting server")
	if err := engine.Run(config.Server.Socket); err != nil {
		panic(err)
	}
	base.Logger.Info("Server stopped")
}

func main() {
	defer processPanic()

	ctx := context.TODO()
	schemaValidator := base.CreateValidator()
	config, err := base.LoadConfiguration(base.ConfigFile)
	if err != nil {
		processError(err)
	}

	setLogger(config)

	mongoClient := createMongoClient(config, &ctx)

	checkMongoConnection(mongoClient, &ctx)

	usersCollection := mongoClient.Database(
		config.MongoDB.Database,
	).Collection(string(base.Users))
	filesCollection := mongoClient.Database(
		config.MongoDB.Database,
	).Collection(string(base.Files))
	filesMetadataCollection := mongoClient.Database(
		config.MongoDB.Database,
	).Collection(string(base.FilesMetadata))

	authService := &services.AuthorizationService{JwtConfig: &config.Server.JwtConfig}
	userService := &services.UserService{Context: &ctx, Collection: usersCollection}
	filesService := &services.FilesService{Context: &ctx, Collection: filesCollection}
	filesMetadataService := &services.FilesMetadataService{
		Context: &ctx, Collection: filesMetadataCollection,
	}

	authController := controllers.AuthorizationController{AuthService: authService}
	tokenController := controllers.TokenController{
		SchemaValidator: schemaValidator,
		AuthService:     authService,
		UserService:     userService,
	}
	userController := controllers.UserController{
		Service:         userService,
		SchemaValidator: schemaValidator,
	}
	filesController := controllers.FilesController{
		FilesService:         filesService,
		FilesMetadataService: filesMetadataService,
		FilesExpConfig:       &config.FilesExpConfig,
		SchemaValidator:      schemaValidator,
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.NoRoute(api.NoRouteHandler)
	router.NoMethod(api.NoMethodHandler)
	router.Use(api.LogsHandler)
	router.Use(api.ErrorHandler)
	router.Use(api.CORSHandler)

	applicationGroup := router.Group(config.Server.BasePath)
	v1 := applicationGroup.Group("/v1")

	v1.GET("/health", controllers.CheckHealth)
	v1.POST("/login", tokenController.SignIn)

	usersGroup := v1.Group("/users")
	usersGroup.POST("", userController.SignUpUser)

	filesGroup := v1.Group("/files")
	filesGroup.GET(
		fmt.Sprintf("/:%s", base.FileIdPathParam),
		filesController.DownloadFile,
	)

	withAuthUsersGroup := v1.Group("/users").Use(authController.Authorize)
	withAuthUsersGroup.GET("/me", userController.GetUser)

	withAuthFilesGroup := v1.Group("/files").Use(authController.Authorize)
	withAuthFilesGroup.POST("", filesController.UploadFile)
	withAuthFilesGroup.GET("", filesController.GetFileMetadataList)

	configureSwagger(applicationGroup, config)

	runServer(router, config)
}
