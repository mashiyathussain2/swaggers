//go:generate $GOPATH/bin/mockgen -destination=./../mock/mock_media.go -package=mock go-app/app Media

package app

import (
	"bytes"
	"context"
	"fmt"

	"go-app/model"
	"go-app/schema"
	"image/jpeg"
	"image/png"
	"strings"
	"time"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Media contains methods to be implemented to register a content service
type Media interface {
	GenerateVideoUploadToken(*schema.GenerateVideoUploadTokenOpts) (*schema.GenerateVideoUploadTokenResp, error)
	CreateVideoMedia(*schema.CreateVideoOpts) (*schema.CreateVideoResp, error)
	CreateImageMedia(opts *schema.CreateImageMediaOpts) (*schema.CreateImageMediaResp, error)
	DeleteMedia(primitive.ObjectID) (bool, error)
	GetVideoMediaByID(primitive.ObjectID) (*schema.GetMediaResp, error)

	GetImageMediaByID(primitive.ObjectID) (*schema.GetMediaResp, error)
}

// MediaImpl implements Content service methods
type MediaImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

//MediaImplOpts contains args required to create a new instance of ContenImpl
type MediaImplOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitMedia returns ContentImpl instance
func InitMedia(opts *MediaImplOpts) Media {
	c := MediaImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
	return &c
}

// GenerateVideoUploadToken generates an aws s3 presigned upload url to upload video to specified bucket s3 bucket
func (mi *MediaImpl) GenerateVideoUploadToken(opts *schema.GenerateVideoUploadTokenOpts) (*schema.GenerateVideoUploadTokenResp, error) {
	fType, err := FileTypeFromFileName(opts.FileName)
	if err != nil {
		return nil, err
	}
	// sending request to aws s3 bucket to return an upload url with passed filename and filetype
	token, err := mi.generateS3UploadToken(opts.FileName, fType)
	if err != nil {
		return nil, err
	}
	return &schema.GenerateVideoUploadTokenResp{Token: token}, nil
}

func (mi *MediaImpl) generateS3UploadToken(videoID, videoType string) (string, error) {
	url, err := mi.App.S3.GetPutObjectRequestURL(&s3.PutObjectInput{
		Bucket: aws.String(mi.App.Config.S3Config.VideoUploadBucket),
		Key:    aws.String(fmt.Sprintf("%s/%s", mi.App.Config.S3Config.VideoUploadPath, videoID+videoType)),
	})
	if err != nil {
		mi.Logger.Err(err).Msg("failed to generate presigned upload url")
		return "", errors.Wrap(err, "failed to generate upload token")
	}
	return url, nil
}

// CreateVideoMedia creates a new video media object in video collection
func (mi *MediaImpl) CreateVideoMedia(opts *schema.CreateVideoOpts) (*schema.CreateVideoResp, error) {
	v := model.Video{
		Type:      model.VideoType,
		GUID:      opts.GUID,
		FileName:  opts.FileName,
		SRCBucket: opts.SRCBucket,
		Dimensions: &model.Dimensions{
			Height: opts.SRCHeight,
			Width:  opts.SRCWidth,
		},
		IsPortrait:       opts.IsPortrait,
		CloudfrontURL:    opts.CloudFrontURL,
		Duration:         opts.Duration,
		Framerate:        opts.Framerate,
		PlaybackBucket:   opts.PlaybackBucket,
		PlaybackURL:      opts.PlaybackURL,
		ThumbnailBuckets: opts.ThumbnailBuckets,
		ThumbnailURLS:    opts.ThumbnailURLS,
		ProcessedAt:      opts.ProcessedAt,
		CreatedAt:        time.Now().UTC(),
	}

	res, err := mi.DB.Collection(model.MediaColl).InsertOne(context.TODO(), v)
	if err != nil {
		mi.Logger.Err(err).Interface("media_info", opts).Msg("failed to create video media")
		return nil, errors.Wrapf(err, "failed to create video media")
	}

	return &schema.CreateVideoResp{
		ID:               res.InsertedID.(primitive.ObjectID),
		GUID:             v.GUID,
		FileName:         v.FileName,
		SRCBucket:        v.SRCBucket,
		Dimensions:       v.Dimensions,
		CloudfrontURL:    v.CloudfrontURL,
		IsPortrait:       v.IsPortrait,
		Duration:         v.Duration,
		Framerate:        v.Framerate,
		PlaybackBucket:   v.PlaybackBucket,
		PlaybackURL:      v.PlaybackURL,
		ThumbnailBuckets: v.ThumbnailBuckets,
		ThumbnailURLS:    v.ThumbnailURLS,
		ProcessedAt:      v.ProcessedAt,
		CreatedAt:        v.CreatedAt,
	}, nil
}

// DeleteMedia delete media document from collection with session
func (mi *MediaImpl) DeleteMedia(id primitive.ObjectID) (bool, error) {
	res, err := mi.DB.Collection(model.MediaColl).DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return false, errors.Wrapf(err, "failed to delete content with id:%s", id.Hex())
	}
	if res.DeletedCount == 0 {
		return false, errors.Errorf("media with id:%s not found", id.Hex())
	}
	return true, nil
}

// GetVideoMediaByID returns video media document with matching id
func (mi *MediaImpl) GetVideoMediaByID(id primitive.ObjectID) (*schema.GetMediaResp, error) {
	var resp schema.GetMediaResp
	if err := mi.DB.Collection(model.MediaColl).FindOne(context.TODO(), bson.M{"_id": id}).Decode(&resp); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return nil, errors.Errorf("media with id:%s not found", id.Hex())
		}
		return nil, errors.Wrapf(err, "query failed to find media with id:%s", id.Hex())
	}
	return &resp, nil
}

// GetImageMediaByID returns image media document with matching id
func (mi *MediaImpl) GetImageMediaByID(id primitive.ObjectID) (*schema.GetMediaResp, error) {
	var resp schema.GetMediaResp
	if err := mi.DB.Collection(model.MediaColl).FindOne(context.TODO(), bson.M{"_id": id}).Decode(&resp); err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return nil, errors.Errorf("media with id:%s not found", id.Hex())
		}
		return nil, errors.Wrapf(err, "query failed to find media with id:%s", id.Hex())
	}
	return &resp, nil
}

// CreateImageMedia takes a base64 source as string, decode, uploads to aws and stores the reference in Image collection
func (mi *MediaImpl) CreateImageMedia(opts *schema.CreateImageMediaOpts) (*schema.CreateImageMediaResp, error) {
	img := IMG{}
	if err := img.DecodeBase64StrToIMG(opts.Base64SRC); err != nil {
		return nil, err
	}
	i := model.Image{
		FileName: strings.ToLower(uuid.NewV1().String()[:4] + opts.FileName),
		FileType: img.Type,
		Dimensions: &model.Dimensions{
			Height: uint(img.Conf.Height),
			Width:  uint(img.Conf.Width),
		},
		CreatedAt: time.Now().UTC(),
	}

	var buf bytes.Buffer
	switch img.Type {
	case "image/png":
		png.Encode(&buf, *img.Img)
	case "image/jpeg":
		jpeg.Encode(&buf, *img.Img, nil)
	case "image/jpg":
		jpeg.Encode(&buf, *img.Img, nil)
	case "default":
		return nil, errors.New("invalid image file type")
	}

	params := s3.PutObjectInput{
		Body:   bytes.NewReader(buf.Bytes()),
		Bucket: aws.String(mi.App.Config.S3Config.ImageUploadBucket),
		Key:    aws.String("/assets/catalog/img/" + i.FileName),
	}

	_, err := mi.App.S3.PutObject(&params)
	if err != nil {
		return nil, errors.Wrap(err, "failed to upload image to cdn")
	}

	i.SRCBucketURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com%s", *params.Bucket, mi.App.Config.S3Config.Region, *params.Key)

	res, err := mi.DB.Collection(model.MediaColl).InsertOne(context.TODO(), i)
	if err != nil {
		mi.Logger.Err(err).Msg("failed to generate image media")
		return nil, errors.Wrap(err, "failed to generate image media")
	}

	resp := schema.CreateImageMediaResp{
		ID:         res.InsertedID.(primitive.ObjectID),
		FileType:   i.FileType,
		FileName:   i.FileName,
		Dimensions: i.Dimensions,
		URL:        i.SRCBucketURL,
	}

	return &resp, nil
}
