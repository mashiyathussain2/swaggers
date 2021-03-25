//go:generate $GOPATH/bin/mockgen -destination=./../mock/mock_content.go -package=mock go-app/app Content

package app

import (
	"context"
	"go-app/model"
	"go-app/schema"
	"time"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Media contains methods to be implemented to register a content service
type Media interface {
	generateS3UploadToken(videoID, videoType string) (string, error)
	GenerateVideoUploadToken(*schema.GenerateVideoUploadTokenOpts) (*schema.GenerateVideoUploadTokenResp, error)
}

// MediaImpl implements Content service methods
type MediaImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
	S3     S3
}

//MediaImplOpts contains args required to create a new instance of ContenImpl
type MediaImplOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
	S3     S3
}

// InitContent returns ContentImpl instance
func InitContent(opts *MediaImplOpts) Media {
	c := MediaImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
		S3:     opts.S3,
	}
	return &c
}

// GenerateVideoUploadToken generates an aws s3 presigned upload url to upload video to specified bucket s3 bucket
func (ci *MediaImpl) GenerateVideoUploadToken(opts *schema.GenerateVideoUploadTokenOpts) (*schema.GenerateVideoUploadTokenResp, error) {
	c := model.Video{
		CreatedAt: time.Now(),
	}
	res, err := ci.DB.Collection(model.VideoContentColl).InsertOne(context.TODO(), c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create video content")
	}

	c.ID = res.InsertedID.(primitive.ObjectID)
	fType, err := FileTypeFromFileName(opts.FileName)
	if err != nil {
		return nil, err
	}
	token, err := ci.generateS3UploadToken(c.ID.Hex(), fType)
	if err != nil {
		return nil, err
	}
	return &schema.GenerateVideoUploadTokenResp{
		ID:    c.ID,
		Token: token,
	}, nil
}

func (ci *MediaImpl) generateS3UploadToken(videoID, videoType string) (string, error) {
	url, err := ci.S3.GetPutObjectRequestURL(&s3.PutObjectInput{
		Bucket: aws.String(ci.App.Config.S3Config.S3VideoUploadBucket),
		Key:    aws.String(videoID + videoType),
	})
	if err != nil {
		ci.Logger.Err(err).Msg("failed to generate presigned upload url")
		return "", errors.Wrap(err, "failed to generate upload token")
	}
	return url, nil
}
