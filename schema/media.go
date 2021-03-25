package schema

import "go.mongodb.org/mongo-driver/bson/primitive"

// GenerateVideoUploadTokenOpts contains fields and validatation for generating a token for uploading a new video directly to s3
type GenerateVideoUploadTokenOpts struct {
	FileName string `json:"file_name" validate:"required"`
}

// GenerateVideoUploadTokenResp contains fields to returned when new video upload token is generated
type GenerateVideoUploadTokenResp struct {
	ID    primitive.ObjectID `json:"id"`
	Token string             `json:"token"`
}
