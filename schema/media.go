package schema

import (
	"go-app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GenerateVideoUploadTokenOpts contains fields and validatation for generating a token for uploading a new video directly to s3
type GenerateVideoUploadTokenOpts struct {
	FileName string `json:"file_name" validate:"required"`
}

// GenerateVideoUploadTokenResp contains fields to returned when new video upload token is generated
type GenerateVideoUploadTokenResp struct {
	Token string `json:"token"`
}

// CreateVideoOpts contains fields and validation required to create a new video document in DB
type CreateVideoOpts struct {
	GUID             string    `json:"guid"`
	FileName         string    `json:"filename"`
	SRCBucket        string    `json:"srcBucket"`
	DestBucket       string    `json:"destBucket"`
	SRCWidth         uint      `json:"srcWidth"`
	SRCHeight        uint      `json:"srcHeight"`
	PlaybackURL      string    `json:"hlsUrl"`
	PlaybackBucket   string    `json:"hlsPlaylist"`
	ThumbnailURLS    []string  `json:"thumbNailsUrls"`
	ThumbnailBuckets []string  `json:"thumbNails"`
	IsPortrait       bool      `json:"isPortrait"`
	Duration         float32   `json:"duration"`
	Framerate        float32   `json:"framerate"`
	CloudFrontURL    string    `json:"cloudFront"`
	ProcessedAt      time.Time `json:"endTime"`
}

// CreateVideoResp contains fields to returned in response to create video
type CreateVideoResp struct {
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	// GUID: to reference task in aws media processing
	GUID string `json:"guid,omitempty" bson:"guid,omitempty"`

	SRCBucket     string `json:"src_bucket,omitempty" bson:"src_bucket,omitempty"`
	FileName      string `json:"filename,omitempty" bson:"filename,omitempty"`
	DestBucket    string `json:"dest_bucket,omitempty" bson:"dest_bucket,omitempty"`
	CloudfrontURL string `json:"cloudfront_url,omitempty" bson:"cloudfront_url,omitempty"`

	IsPortrait bool              `json:"is_portrait,omitempty" bson:"is_portrait,omitempty"`
	Dimensions *model.Dimensions `json:"dimensions,omitempty" bson:"dimensions,omitempty"`
	Duration   float32           `json:"duration,omitempty" bson:"duration,omitempty"`
	Framerate  float32           `json:"framerate,omitempty" bson:"framerate,omitempty"`

	PlaybackBucket string `json:"hls_playback_bucket,omitempty" bson:"hls_playback_bucket,omitempty"`
	PlaybackURL    string `json:"hls_playback_url,omitempty" bson:"hls_playback_url,omitempty"`

	ThumbnailBuckets []string `json:"thumbnail_bucket,omitempty" bson:"thumbnail_bucket,omitempty"`
	ThumbnailURLS    []string `json:"thumbnail_url,omitempty" bson:"thumbnail_url,omitempty"`

	CreatedAt   time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ProcessedAt time.Time `json:"processed_at,omitempty" bson:"processed_at,omitempty"`
}

// GetMediaResp contains fields to returned in GetMedia response
type GetMediaResp struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FileName      string             `json:"filename,omitempty" bson:"filename,omitempty"`
	CreatedAt     time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	SRCBucket     string             `json:"src_bucket,omitempty" bson:"src_bucket,omitempty"`
	CloudfrontURL string             `json:"cloudfront_url,omitempty" bson:"cloudfront_url,omitempty"`
	Dimensions    *model.Dimensions  `json:"dimensions,omitempty" bson:"dimensions,omitempty"`

	// Video
	IsPortrait    bool     `json:"is_portrait,omitempty" bson:"is_portrait,omitempty"`
	Duration      float32  `json:"duration,omitempty" bson:"duration,omitempty"`
	Framerate     float32  `json:"framerate,omitempty" bson:"framerate,omitempty"`
	PlaybackURL   string   `json:"hls_playback_url,omitempty" bson:"hls_playback_url,omitempty"`
	ThumbnailURLS []string `json:"thumbnail_url,omitempty" bson:"thumbnail_url,omitempty"`

	// Image
	FileType string `json:"file_type,omitempty" bson:"file_type,omitempty"`
	URL      string `json:"url,omitempty" bson:"url,omitempty"`
}

// CreateImageMediaOpts contains fields and validations required to create image media
type CreateImageMediaOpts struct {
	FileName  string `json:"file_name" validate:"required"`
	Base64SRC string `json:"base64_src" validate:"required"`
}

// CreateImageMediaResp contains fields to be returned for image upload
type CreateImageMediaResp struct {
	ID         primitive.ObjectID `json:"id"`
	FileName   string             `json:"file_name"`
	FileType   string             `json:"file_type"`
	Dimensions *model.Dimensions  `json:"dimensions"`
	URL        string             `json:"url"`
}
