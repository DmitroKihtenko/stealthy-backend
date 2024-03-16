package services

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	mongoMock "github.com/sv-tools/mongoifc/mocks/mockery"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"stealthy-backend/api"
	"stealthy-backend/base"
	"stealthy-backend/tests"
	"testing"
)

func TestCheckFileExists(t *testing.T) {
	opts := options.FindOne().SetProjection(bson.D{{Key: "identifier", Value: 1}})
	dbContext := context.TODO()
	fileData := tests.FileDataFactory.Build()
	emptyFileData := api.FileData{}

	collectionMock := new(mongoMock.Collection)
	resultMock := new(mongoMock.SingleResult)
	resultMock.On("Decode", &emptyFileData).Return(nil)
	collectionMock.On("FindOne", dbContext, bson.D{
		primitive.E{Key: "identifier", Value: fileData.Identifier},
	}, opts).Return(resultMock)

	service := FilesService{Context: &dbContext, Collection: collectionMock}
	result, err := service.CheckFileDataExists(fileData.Identifier)

	assert.Equal(t, true, result)
	assert.Nil(t, err)
}

func TestCheckFileDoesNotExist(t *testing.T) {
	opts := options.FindOne().SetProjection(bson.D{{Key: "identifier", Value: 1}})
	dbContext := context.TODO()
	fileData := tests.FileDataFactory.Build()
	emptyFileData := api.FileData{}

	collectionMock := new(mongoMock.Collection)
	resultMock := new(mongoMock.SingleResult)
	resultMock.On("Decode", &emptyFileData).Return(mongo.ErrNoDocuments)
	collectionMock.On("FindOne", dbContext, bson.D{
		primitive.E{Key: "identifier", Value: fileData.Identifier},
	}, opts).Return(resultMock)

	service := FilesService{Context: &dbContext, Collection: collectionMock}
	result, err := service.CheckFileDataExists(fileData.Identifier)

	assert.Equal(t, false, result)
	assert.Nil(t, err)
}

func TestAddFile(t *testing.T) {
	opts := options.FindOne().SetProjection(bson.D{{Key: "identifier", Value: 1}})
	dbContext := context.TODO()
	fileDataToAdd := tests.FileDataFactory.Build()
	emptyFileData := api.FileData{}
	expectedResult := api.AddFileResponse{Identifier: fileDataToAdd.Identifier}

	collectionMock := new(mongoMock.Collection)
	resultMock := new(mongoMock.SingleResult)
	resultMock.On("Decode", &emptyFileData).Return(mongo.ErrNoDocuments)
	collectionMock.On("FindOne", dbContext, bson.D{
		primitive.E{Key: "identifier", Value: fileDataToAdd.Identifier},
	}, opts).Return(resultMock)
	collectionMock.On("InsertOne", dbContext, &fileDataToAdd).Return(
		nil, nil,
	)

	service := FilesService{Context: &dbContext, Collection: collectionMock}
	result, err := service.AddFile(&fileDataToAdd)

	assert.Equal(t, expectedResult, *result)
	assert.Nil(t, err)
}

func TestAddFileAlreadyExist(t *testing.T) {
	opts := options.FindOne().SetProjection(bson.D{{Key: "identifier", Value: 1}})
	dbContext := context.TODO()
	fileDataToAdd := tests.FileDataFactory.Build()
	expectedError := base.ServiceError{
		Summary: fmt.Sprintf("File '%s' already exist", fileDataToAdd.Identifier),
		Status:  http.StatusBadRequest,
	}

	collectionMock := new(mongoMock.Collection)
	resultMock := new(mongoMock.SingleResult)
	resultMock.On("Decode", &api.FileData{}).Return(nil)
	collectionMock.On("FindOne", dbContext, bson.D{
		primitive.E{Key: "identifier", Value: fileDataToAdd.Identifier},
	}, opts).Return(resultMock)
	collectionMock.On("InsertOne", dbContext, &fileDataToAdd).Return(
		nil, nil,
	)

	service := FilesService{Context: &dbContext, Collection: collectionMock}
	result, err := service.AddFile(&fileDataToAdd)

	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestGetFile(t *testing.T) {
	dbContext := context.TODO()
	fileData := tests.FileDataFactory.Build()
	emptyFileData := api.FileData{}

	collectionMock := new(mongoMock.Collection)
	resultMock := new(mongoMock.SingleResult)
	resultMock.On("Decode", &emptyFileData).Return(func(
		v interface{},
	) error {
		if v != nil {
			copier.Copy(v, &fileData)
		}
		return nil
	})
	collectionMock.On("FindOne", dbContext, bson.D{
		primitive.E{Key: "identifier", Value: fileData.Identifier},
	}).Return(resultMock)

	service := FilesService{Context: &dbContext, Collection: collectionMock}
	result, err := service.GetFile(fileData.Identifier)

	assert.Equal(t, &fileData, result)
	assert.Nil(t, err)
}

func TestGetFileNotFound(t *testing.T) {
	dbContext := context.TODO()
	fileData := tests.FileDataFactory.Build()
	emptyFileData := api.FileData{}

	collectionMock := new(mongoMock.Collection)
	resultMock := new(mongoMock.SingleResult)
	resultMock.On("Decode", &emptyFileData).Return(mongo.ErrNoDocuments)
	collectionMock.On("FindOne", dbContext, bson.D{
		primitive.E{Key: "identifier", Value: fileData.Identifier},
	}).Return(resultMock)

	service := FilesService{Context: &dbContext, Collection: collectionMock}
	result, err := service.GetFile(fileData.Identifier)

	assert.Nil(t, result)
	assert.Equal(t, base.ServiceError{
		Summary: fmt.Sprintf("'%s' not found", fileData.Identifier),
		Status:  http.StatusNotFound,
	}, err)
}
