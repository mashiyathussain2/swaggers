package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// list of collection name
const (
	LiveColl string = "live"
)

// StreamStatus represents status of the stream
type StreamStatus struct {
	Name      string    `json:"name,omitempty" bson:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// list of supported stream status
const (
	ActiveStatus  string = "active"
	DiscardStatus string = "discard"
	EndStatus     string = "end"
)

// Live contains hypd live stream data.
type Live struct {
	ID            primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string               `json:"name,omitempty" bson:"name,omitempty"`
	Slug          string               `json:"slug,omitempty" bson:"slug,omitempty"`
	InfluencerIDs []primitive.ObjectID `json:"influencer_ids" bson:"influencer_ids,omitempty"`
	Status        *StreamStatus        `json:"status,omitempty" bson:"status,omitempty"`
	StatusHistory []StreamStatus       `json:"status_history,omitempty" bson:"status_history,omitempty"`

	LikeCount    uint                 `json:"like_count" bson:"like_count"`
	LikeIDs      []primitive.ObjectID `json:"like_ids,omitempty" bson:"like_ids,omitempty"`
	LikedBy      []primitive.ObjectID `json:"liked_by,omitempty" bson:"liked_by,omitempty"`
	CommentCount uint                 `json:"comment_count" bson:"comment_count"`

	CatalogIDs     []primitive.ObjectID `json:"catalog_ids,omitempty" bson:"catalog_ids,omitempty"`
	FeaturedImage  *IMG                 `json:"featured_image,omitempty" bson:"featured_image,omitempty"`
	StreamEndImage *IMG                 `json:"stream_end_image,omitempty" bson:"stream_end_image,omitempty"`
	IVS            *IVS                 `json:"ivs,omitempty" bson:"ivs,omitempty"`
	ScheduledAt    time.Time            `json:"scheduled_at,omitempty" bson:"scheduled_at,omitempty"`
	CreatedAt      time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

// IVSChannel contains channel specific configuration
type IVSChannel struct {
	ARN                   string `json:"arn,omitempty" bson:"arn,omitempty"`
	Name                  string `json:"name,omitempty" bson:"name,omitempty"`
	Type                  string `json:"type,omitempty" bson:"type,omitempty"`
	LatencyMode           string `json:"latency_mode,omitempty" bson:"latency_mode,omitempty"`
	PlaybackAuthorization bool   `json:"playback_authorization,omitempty" bson:"playback_authorization,omitempty"`
}

// IVSIngest contains channel video ingestion specific configuration
type IVSIngest struct {
	IngestURL string `json:"server_url,omitempty" bson:"server_url,omitempty"`
	StreamKey string `json:"stream_key,omitempty" bson:"stream_key,omitempty"`
}

// IVSPlayback contains IVS playback specific configuration
type IVSPlayback struct {
	PlaybackURL string `json:"playback_url,omitempty" bson:"playback_url,omitempty"`
}

// IVS contains aws IVS specific configuration
type IVS struct {
	Channel   *IVSChannel  `json:"channel,omitempty" bson:"channel,omitempty"`
	Ingestion *IVSIngest   `json:"ingestion,omitempty" bson:"ingestion,omitempty"`
	Playback  *IVSPlayback `json:"playback,omitempty" bson:"playback,omitempty"`
}
