package model

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"syreclabs.com/go/faker"
)

// GetRandomContent returns Content model with random data
func GetRandomContent() *Content {
	c := Content{
		Caption: faker.Lorem().Sentence(faker.RandomInt(20, 200)),
		Label: &Label{
			Interests: faker.Lorem().Words(faker.RandomInt(1, 10)),
		},
		CreatedAt:   time.Now().UTC(),
		ProcessedAt: time.Now().UTC(),
		IsProcessed: true,
	}
	for i := 0; i < faker.RandomInt(1, 10); i++ {
		c.Hashtags = append(c.Hashtags, faker.Letterify(faker.RandomChoice([]string{"#???????", "#??????", "#???", "#?????????"})))
	}
	for i := 0; i < faker.RandomInt(1, 2); i++ {
		x := []string{"M", "F", "O"}
		c.Label.Genders = append(c.Label.Genders, x[i])
	}
	for i := 0; i < faker.RandomInt(1, 4); i++ {
		x := []string{"18-22", "23-30", "31-40", "40-50", "50-"}
		c.Label.AgeGroups = append(c.Label.AgeGroups, x[i])
	}
	for i := 0; i < faker.RandomInt(1, 10); i++ {
		c.InfluencerIDs = append(c.InfluencerIDs, primitive.NewObjectIDFromTimestamp(time.Now()))
	}
	for i := 0; i < faker.RandomInt(1, 4); i++ {
		c.BrandIDs = append(c.BrandIDs, primitive.NewObjectIDFromTimestamp(time.Now()))
	}
	for i := 0; i < faker.RandomInt(1, 10); i++ {
		c.CatalogIDs = append(c.CatalogIDs, primitive.NewObjectIDFromTimestamp(time.Now()))
	}
	return &c
}

// GetRandomVideoMedia returns video media populated with random data
func GetRandomVideoMedia() *Video {
	v := &Video{
		Type:       VideoType,
		GUID:       faker.Letterify("????????-????-????-????-????????????"),
		FileName:   fmt.Sprintf("%s%s", primitive.NewObjectIDFromTimestamp(time.Now()).Hex(), faker.RandomChoice([]string{".mp4", ".mov"})),
		SRCBucket:  faker.Letterify("hypd-vod-source-??????????"),
		DestBucket: faker.Letterify("hypd-vod-dest-??????????"),
		Dimensions: &Dimensions{
			Width:  uint(faker.RandomInt(360, 1080)),
			Height: uint(faker.RandomInt(360, 1080)),
		},
		Duration:         float32(faker.RandomInt(1, 1000)),
		Framerate:        float32(faker.RandomInt(20, 60)),
		IsPortrait:       false,
		CloudfrontURL:    faker.Letterify("??????????.cloudfront.net"),
		ProcessedAt:      time.Now().UTC(),
		PlaybackBucket:   "s3://hypd-vod-destination-r955ikuyz5i4/bee228c2-2f09-4c92-8c60-832b107003d5/hls/603236a958b8f4136b6e44a1.m3u8",
		PlaybackURL:      "https://d26egzot5z0rj9.cloudfront.net/bee228c2-2f09-4c92-8c60-832b107003d5/hls/603236a958b8f4136b6e44a1.m3u8",
		ThumbnailBuckets: []string{"s3://hypd-vod-destination-r955ikuyz5i4/bee228c2-2f09-4c92-8c60-832b107003d5/thumbnails/603236a958b8f4136b6e44a1_thumb.0000005.jpg"},
		ThumbnailURLS:    []string{"https://d26egzot5z0rj9.cloudfront.net/bee228c2-2f09-4c92-8c60-832b107003d5/thumbnails/603236a958b8f4136b6e44a1_thumb.0000005.jpg"},
	}
	return v
}

// GetRandomLive returns live model populated with random data
func GetRandomLive() *Live {
	name := faker.Name().Name()
	now := time.Now().UTC()
	l := &Live{
		ID:            primitive.NewObjectIDFromTimestamp(time.Now()),
		Name:          name,
		Slug:          faker.Internet().Slug(),
		InfluencerIDs: []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now())},
		CatalogIDs:    []primitive.ObjectID{primitive.NewObjectIDFromTimestamp(time.Now()), primitive.NewObjectIDFromTimestamp(time.Now())},
		Status: &StreamStatus{
			Name:      ActiveStatus,
			CreatedAt: now,
		},
		StatusHistory: []StreamStatus{
			{
				Name:      ActiveStatus,
				CreatedAt: now,
			},
			{
				Name:      DiscardStatus,
				CreatedAt: now.Add(-36 * time.Hour),
			},
		},
		FeaturedImage: &IMG{
			SRC:    "https://deepakacademy.files.wordpress.com/2020/07/carryminatiyu.jpg",
			Width:  1280,
			Height: 720,
		},
		StreamEndImage: &IMG{
			SRC:    "https://i2.wp.com/www.movieslantern.com/wp-content/uploads/2019/10/maxresdefault-170.jpg?fit=768%2C432&ssl=1",
			Width:  1280,
			Height: 720,
		},
		IVS: &IVS{
			Channel: &IVSChannel{
				ARN:                   faker.Letterify("???-????-?????"),
				Name:                  name,
				Type:                  "STANDARD",
				LatencyMode:           "LOW",
				PlaybackAuthorization: false,
			},
			Ingestion: &IVSIngest{
				IngestURL: faker.Letterify("rtmp://??.???.??.??:8000"),
				StreamKey: faker.RandomString(20),
			},
			Playback: &IVSPlayback{
				PlaybackURL: faker.Letterify("https://?????.??????.com/????.m3u8"),
			},
		},
		ScheduledAt: time.Now().Add(time.Duration(faker.RandomInt(10, 1000) * int(time.Hour))).UTC(),
		CreatedAt:   time.Now().UTC(),
	}
	return l
}
