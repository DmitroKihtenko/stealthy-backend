package services

import (
	"SharingBackend/api"
	"SharingBackend/base"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func CheckPasswordEquals(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type UserService struct {
	Context    *context.Context
	Collection *mongo.Collection
}

func (service *UserService) checkUserExists(request *api.SignUpRequest) (bool, error) {
	var user api.User
	err := service.Collection.FindOne(*service.Context, bson.D{primitive.E{
		Key: "username", Value: request.Username},
	}).Decode(&user)

	if err == nil {
		return true, nil
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	} else {
		return false, base.NewDatabaseError(err)
	}
}

func (service *UserService) AddUser(request *api.SignUpRequest) (*api.UserResponse, error) {
	exists, err := service.checkUserExists(request)
	if exists && err != nil {
		return nil, err
	} else if exists {
		err := base.ServiceError{
			Summary: fmt.Sprintf("User '%s' already exist", request.Username),
			Status:  http.StatusBadRequest,
		}
		return nil, err
	} else {
		bytes, err := bcrypt.GenerateFromPassword([]byte(request.Password), base.PasswordCost)
		user := api.User{
			Username:     request.Username,
			PasswordHash: string(bytes),
		}
		if err != nil {
			err := base.ServiceError{
				Summary: "Password processing error",
				Detail:  err.Error(),
			}
			return nil, err
		}

		if _, err := service.Collection.InsertOne(*service.Context, &user); err != nil {
			return nil, base.NewDatabaseError(err)
		}
		return &api.UserResponse{Username: request.Username}, nil
	}
}

func (service *UserService) GetUserByUsername(username string) (*api.User, error) {
	var user api.User
	err := service.Collection.FindOne(*service.Context, bson.D{primitive.E{
		Key: "username", Value: username},
	}).Decode(&user)

	if err == nil {
		return &user, nil
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, base.ServiceError{
			Summary: fmt.Sprintf("User '%s' not found", username),
			Status:  http.StatusNotFound,
		}
	} else {
		return nil, base.NewDatabaseError(err)
	}
}

func (service *UserService) GetUserPublicData(username string) (*api.UserResponse, error) {
	user, err := service.GetUserByUsername(username)
	if err != nil {
		return nil, err
	} else {
		return &api.UserResponse{
			Username: user.Username,
		}, nil
	}
}

func (service *UserService) GetUserByCredentials(request *api.SignInRequest) (*api.User, error) {
	user, err := service.GetUserByUsername(request.Username)
	if err != nil {
		return nil, err
	}

	if CheckPasswordEquals(request.Password, user.PasswordHash) {
		return user, nil
	} else {
		return nil, base.ServiceError{
			Summary: fmt.Sprintf("Invalid password"),
			Status:  http.StatusForbidden,
		}
	}
}
