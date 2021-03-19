package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// list of collection name
const (
	MediaColl string = "media"
)

// list of supported media types
const (
	VideoType string = "video"
	ImageType string = "image"
)

// Dimensions contains height and width of video in pixels
type Dimensions struct {
	Height uint `json:"height,omitempty" bson:"height,omitempty"`
	Width  uint `json:"width,omitempty" bson:"width,omitempty"`
}

// Video contains video content data such as type, url, source-bucket, content meta etc
type Video struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Type string             `json:"type,omitempty" bson:"_type,omitempty"`
	// GUID: to reference task in aws media processing
	GUID string `json:"guid,omitempty" bson:"guid,omitempty"`

	SRCBucket     string `json:"src_bucket,omitempty" bson:"src_bucket,omitempty"`
	FileName      string `json:"filename,omitempty" bson:"filename,omitempty"`
	DestBucket    string `json:"dest_bucket,omitempty" bson:"dest_bucket,omitempty"`
	CloudfrontURL string `json:"cloudfront_url,omitempty" bson:"cloudfront_url,omitempty"`

	IsPortrait bool        `json:"is_portrait,omitempty" bson:"is_portrait,omitempty"`
	Dimensions *Dimensions `json:"dimensions,omitempty" bson:"dimensions,omitempty"`
	Duration   float32     `json:"duration,omitempty" bson:"duration,omitempty"`
	Framerate  uint        `json:"framerate,omitempty" bson:"framerate,omitempty"`

	PlaybackBucket string `json:"hls_playback_bucket,omitempty" bson:"hls_playback_bucket,omitempty"`
	PlaybackURL    string `json:"hls_playback_url,omitempty" bson:"hls_playback_url,omitempty"`

	ThumbnailBuckets []string `json:"thumbnail_bucket,omitempty" bson:"thumbnail_bucket,omitempty"`
	ThumbnailURLS    []string `json:"thumbnail_url,omitempty" bson:"thumbnail_url,omitempty"`

	CreatedAt   time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ProcessedAt time.Time `json:"processed_at,omitempty" bson:"processed_at,omitempty"`
}

// Image contains image content data such as url, meta etc
type Image struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FileName   string             `json:"file_name,omitempty" bson:"file_name,omitempty"`
	FileType   string             `json:"file_type,omitempty" bson:"file_type,omitempty"`
	Dimensions *Dimensions        `json:"dimensions,omitempty" bson:"dimensions,omitempty"`

	SRCBucket string `json:"src_bucket,omitempty" bson:"src_bucket,omitempty"`
	// To access image from s3 bucket
	SRCBucketURL  string `json:"src_bucket_url,omitempty" bson:"src_bucket_url,omitempty"`
	CloudfrontURL string `json:"cloudfront_url,omitempty" bson:"cloudfront_url,omitempty"`
	// URL = CloudfrontURL + SRCBucket (used by app); to access image from cloudfront
	URL string `json:"url,omitempty" bson:"url,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
