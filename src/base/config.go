package base

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type MongoDBConfig struct {
	URL            string `yaml:"url" validate:"required"`
	Database       string `yaml:"database" validate:"required"`
	SecondsTimeout int    `yaml:"secondsTimeout" validate:"required,gt=0"`
}

type JwtConfig struct {
	Secret       string `yaml:"secret" validate:"required"`
	DaysLifespan int    `yaml:"DaysLifespan" validate:"required,gt=0"`
}

type ServerConfig struct {
	Socket                 string    `yaml:"socket" validate:"required,unix_addr"`
	BasePath               string    `yaml:"basePath"`
	OpenapiBasePath        string    `yaml:"openapiBasePath"`
	JwtConfig              JwtConfig `yaml:"jwtConfig" validate:"required"`
	PaginationDefaultLimit int64     `yaml:"paginationDefaultLimit" validate:"required,gt=1"`
}

type FilesExpirationConfig struct {
	MinutesLifetimeDefault uint64 `yaml:"minutesLifetimeDefault" validate:"required,gt=0"`
}

type LogConfig struct {
	Level   string `yaml:"level" validate:"required,oneof=fatal error warn warning info debug trace"`
	AppName string `yaml:"appName" validate:"required"`
}

type BackendConfig struct {
	MongoDB        MongoDBConfig         `yaml:"mongoDB"`
	Server         ServerConfig          `yaml:"server"`
	FilesExpConfig FilesExpirationConfig `yaml:"filesExpConfig"`
	Logs           LogConfig             `yaml:"logs"`
}

func LoadConfiguration(file string) (*BackendConfig, error) {
	Logger.WithFields(logrus.Fields{"filename": file}).Info(
		"Loading configuration",
	)

	cfg := &BackendConfig{}
	cfg.SetDefaults()
	if err := cfg.loadFromFile(file); err != nil {
		return nil, err
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *BackendConfig) SetDefaults() {
	cfg.MongoDB.Database = "sharing-backend"
	cfg.MongoDB.URL = "mongodb://root:secret@mongodb:27017"

	cfg.Server.Socket = "localhost:8000"
	cfg.Server.BasePath = "/backend"
	cfg.Server.OpenapiBasePath = "/swagger"
	cfg.Server.PaginationDefaultLimit = 20

	cfg.Server.JwtConfig.DaysLifespan = 3

	cfg.FilesExpConfig.MinutesLifetimeDefault = 1

	cfg.Logs.Level = logrus.DebugLevel.String()
	cfg.Logs.AppName = "sharing-backend"
}

func (cfg *BackendConfig) loadFromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("config file '%s' open error. %s", file, err.Error())
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return fmt.Errorf(
			"config file '%s' reading error, invalid format. %s",
			file,
			err.Error(),
		)
	}
	return nil
}

func (cfg *BackendConfig) validate() error {
	validatorObj := validator.New()
	if err := validatorObj.Struct(cfg); err != nil {
		return WrapValidationErrors(err)
	}
	return nil
}
