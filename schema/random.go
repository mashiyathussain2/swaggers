package schema

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	faker "syreclabs.com/go/faker"
)

// GetRandomCreatePebbleOpts returns CreatePebbleOpts populated with random data
func GetRandomCreatePebbleOpts() *CreatePebbleOpts {
	s := &CreatePebbleOpts{
		Caption:  faker.Lorem().Sentence(faker.RandomInt(20, 200)),
		FileName: fmt.Sprintf("%s%s", faker.Lorem().Word(), faker.RandomChoice([]string{".mp4", ".mov", ".mpg", ".m4v", ".m2ts"})),
		Label: &LabelOpts{
			Interests: faker.Lorem().Words(faker.RandomInt(1, 10)),
		},
	}

	for i := 0; i < faker.RandomInt(1, 10); i++ {
		s.Label.Gender = append(s.Label.Gender, faker.RandomChoice([]string{"M", "F", "O"}))
	}
	for i := 0; i < faker.RandomInt(1, 10); i++ {
		s.Label.AgeGroup = append(s.Label.AgeGroup, faker.RandomChoice([]string{"18-22", "23-30", "31-40", "40-50", "50-"}))
	}
	for i := 0; i < faker.RandomInt(1, 10); i++ {
		s.InfluencerIDs = append(s.InfluencerIDs, primitive.NewObjectIDFromTimestamp(time.Now()))
	}
	for i := 0; i < faker.RandomInt(1, 4); i++ {
		s.BrandIDs = append(s.BrandIDs, primitive.NewObjectIDFromTimestamp(time.Now()))
	}
	for i := 0; i < faker.RandomInt(1, 10); i++ {
		s.CatalogIDs = append(s.CatalogIDs, primitive.NewObjectIDFromTimestamp(time.Now()))
	}
	return s
}

// GetRandomCreateVideoOpts returns CreateVideoOpts populated with random data
func GetRandomCreateVideoOpts() *CreateVideoOpts {
	v := &CreateVideoOpts{
		GUID:             faker.Letterify("????????-????-????-????-????????????"),
		FileName:         fmt.Sprintf("%s%s", primitive.NewObjectIDFromTimestamp(time.Now()).Hex(), faker.RandomChoice([]string{".mp4", ".mov"})),
		SRCBucket:        faker.Letterify("hypd-vod-source-??????????"),
		DestBucket:       faker.Letterify("hypd-vod-dest-??????????"),
		SRCWidth:         uint(faker.RandomInt(360, 1080)),
		SRCHeight:        uint(faker.RandomInt(360, 1080)),
		Duration:         float32(faker.RandomInt(1, 1000)),
		Framerate:        uint(faker.RandomInt(20, 60)),
		IsPortrait:       false,
		CloudFrontURL:    faker.Letterify("??????????.cloudfront.net"),
		ProcessedAt:      time.Now().UTC(),
		PlaybackBucket:   "s3://hypd-vod-destination-r955ikuyz5i4/bee228c2-2f09-4c92-8c60-832b107003d5/hls/603236a958b8f4136b6e44a1.m3u8",
		PlaybackURL:      "https://d26egzot5z0rj9.cloudfront.net/bee228c2-2f09-4c92-8c60-832b107003d5/hls/603236a958b8f4136b6e44a1.m3u8",
		ThumbnailBuckets: []string{"s3://hypd-vod-destination-r955ikuyz5i4/bee228c2-2f09-4c92-8c60-832b107003d5/thumbnails/603236a958b8f4136b6e44a1_thumb.0000005.jpg"},
		ThumbnailURLS:    []string{"https://d26egzot5z0rj9.cloudfront.net/bee228c2-2f09-4c92-8c60-832b107003d5/thumbnails/603236a958b8f4136b6e44a1_thumb.0000005.jpg"},
	}

	if faker.RandomInt(0, 1) == 1 {
		v.IsPortrait = true
	}

	return v
}

// GetRandomCreateCatalogContentOpts returns CreateCatalogContentOpts populated with random data
func GetRandomCreateCatalogContentOpts() *CreateVideoCatalogContentOpts {
	s := &CreateVideoCatalogContentOpts{
		FileName: fmt.Sprintf("%s%s", faker.Lorem().Word(), faker.RandomChoice([]string{".mp4", ".mov", ".mpg", ".m4v", ".m2ts"})),
		Label: &LabelOpts{
			Interests: faker.Lorem().Words(faker.RandomInt(1, 10)),
		},
		BrandID:   primitive.NewObjectIDFromTimestamp(time.Now().UTC()),
		CatalogID: primitive.NewObjectIDFromTimestamp(time.Now().UTC()),
	}

	for i := 0; i < faker.RandomInt(1, 10); i++ {
		s.Label.Gender = append(s.Label.Gender, faker.RandomChoice([]string{"M", "F", "O"}))
	}
	for i := 0; i < faker.RandomInt(1, 10); i++ {
		s.Label.AgeGroup = append(s.Label.AgeGroup, faker.RandomChoice([]string{"18-22", "23-30", "31-40", "40-50", "50-"}))
	}
	return s
}
