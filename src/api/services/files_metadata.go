package services

import (
	"SharingBackend/api"
	"SharingBackend/base"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

type FilesMetadataService struct {
	Context    *context.Context
	Collection *mongo.Collection
}

func (service FilesMetadataService) CheckFileMetadataExists(fileId string) (bool, error) {
	var data api.FileData
	opts := options.FindOne().SetProjection(bson.D{{Key: "identifier", Value: 1}})
	err := service.Collection.FindOne(*service.Context, bson.D{
		primitive.E{Key: "identifier", Value: fileId},
	}, opts).Decode(&data)

	if err == nil {
		return true, nil
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	} else {
		return false, base.NewDatabaseError(err)
	}
}

func (service FilesMetadataService) AddFileMetadata(request *api.FileMetadata) (*api.AddFileResponse, error) {
	exists, err := service.CheckFileMetadataExists(request.Identifier)
	if exists && err != nil {
		return nil, err
	} else if exists {
		err := base.ServiceError{
			Summary: fmt.Sprintf("File '%s' already exist", request.Identifier),
			Status:  http.StatusBadRequest,
		}
		return nil, err
	} else {
		if _, err := service.Collection.InsertOne(*service.Context, &request); err != nil {
			return nil, base.NewDatabaseError(err)
		}
		return &api.AddFileResponse{Identifier: request.Identifier}, nil
	}
}

func (service FilesMetadataService) GetFileMetadataList(
	queryParams *api.PaginationQueryParameters,
	username string,
) (*api.FileMetadataListResponse, error) {
	metadataListResponse := api.FileMetadataListResponse{}

	findOptions := options.Find().
		SetSkip(queryParams.Skip).
		SetLimit(queryParams.Limit).
		SetSort(bson.M{"creation": -1})
	filter := bson.D{
		primitive.E{Key: "username", Value: username},
	}

	total, err := service.Collection.CountDocuments(*service.Context, filter)
	if err != nil {
		err := base.NewDatabaseError(err)
		return nil, err
	}

	metadataListResponse.Total = total
	metadataListResponse.Records = []*api.FileMetadata{}

	cursor, err := service.Collection.Find(*service.Context, filter, findOptions)
	if err != nil {
		err := base.NewDatabaseError(err)
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx *context.Context) {
		err := cursor.Close(*ctx)
		if err != nil {
			base.Logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Warn("Close cursor error")
		}
	}(cursor, service.Context)

	for cursor.Next(context.Background()) {
		var fileMetadata api.FileMetadata
		if err := cursor.Decode(&fileMetadata); err != nil {
			err := base.NewDatabaseError(err)
			return nil, err
		}
		metadataListResponse.Records = append(metadataListResponse.Records, &fileMetadata)
	}

	return &metadataListResponse, nil
}

func (service FilesMetadataService) GetFileMetadata(
	fileId string,
) (*api.FileMetadata, error) {
	var fileData api.FileMetadata
	err := service.Collection.FindOne(*service.Context, bson.D{
		primitive.E{Key: "identifier", Value: fileId},
	}).Decode(&fileData)

	if err == nil {
		return &fileData, nil
	} else if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, base.ServiceError{
			Summary: fmt.Sprintf("File '%s' not found", fileId),
			Status:  http.StatusNotFound,
		}
	} else {
		return nil, base.NewDatabaseError(err)
	}
}
