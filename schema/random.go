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

func GetRandomCreateImageMediaOpts() *CreateImageMediaOpts {
	i := CreateImageMediaOpts{
		FileName:  "hypd.png",
		Base64SRC: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAOEAAADhCAMAAAAJbSJIAAAAe1BMVEUAAAD////8/PwEBATOzs6/v7+ysrLx8fFDQ0OHh4f5+fkICAjExMRnZ2e2trYyMjI6OjpWVlbJycnT09Pt7e3g4OCfn59eXl7R0dGUlJRSUlJ6eno/Pz9qamrm5ubb29upqakhISFLS0uAgIBycnIZGRksLCwdHR2jo6PklQJVAAAEyklEQVR4nO2ZDVPqOhCGm5SCNoBQUPw4ChfRc///L7zZJNq0CeLxgjhnnofRmabbNG93s9m0RQEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAcG7KtMl0ThsxKTtN9te9rEx66TeUcUsZ/Z0UGagpjemNVprej4yMvqPIFKa0mD1XlIH+vUzUFixOLrC9+yda/ux6J+KQjUmfw7F5HjTNoFms2pb1wHI7jfyzamyLtbpvm+6kwdK27BauwdM0u6siDdzJYhDZDG5mL98QpQul7U9dti3KMbpoWx6lQWt1/d5yr4JVazRWPW4fkpsN+zbXN6eXOFBO4ocKi+swoMdCpp71TOMk1yoSUfVHr9Xtr16kJgqt1WNhjInz2lkUPoibLQOfIIoLObYCnyKbRGFdq9G8mz8ThbrW6k4y0LkVloswol3hEujYOrB3VcaH1mQ0Kz70YS0SJ5mF5psVmmKlaj9vXiROVzJ627COO8r40NqMOzdLfSj9LF9POhc/o7AMCcmG5kQOnAvt0HapQt0dv5u6XYU6tqnFRG1OqO+TUWqKrUy7MPMmMnRJrb/jyqfyktab4XC4GTa1E6K7fg4Ka7EZbqqlewa2039+gMKi0l7Wpcus7tGvirKrUEyG4XC79lbLbdRPiNI6HM4ad0ltI+OEfCpKbSaYh8E8FTsJLu+dMvXh9O343sW11vHoOwpt9JqRfwzPp1InfMqHkuwGfsRqt3ZBqmQtjHNg5aZVUChrxFL0yWLQEit0y86t7+u6OCFB4WY1CcxzCi0XPitIcnAjX/TOd3woVfjAPZC9CsNFuttwApxCO/nrXh5PFNrh6ZAutchc9U5Xvsrx81CiN5g/d7rIKxwVJ8T7UGt9UOGqtcm40PtQtfNw6JeOzpM4n0LV92BGYVls3vSJxHl/lQ5ROn6dWX7dT5WfYnY5bzmfwlod9KHVM6u9OBl4lXRUqWg1116uyzxp1XYWH35GYTEVX/uI3hb97UBmb2ENl1exzY+eh5LcX7QftsuY/f1ArvKWVShTeX+7QvFg3bJfocseLhKv0u1AzocSzD9Bodzk7vUqYPIKXWFzVXtvTzObgXRvofQ0vwM+w3p4sKZ5Y+Sdc5k5Vb0Fpo9Q++8pKTc7NY283wt3j3fS/s3fEXeMx1So3Z7dX76uJr8Tm2GoF9yB1fAy8rk56q8sfQ77mQpltNV2blld+NH2xhmqgPe9xfin7C1aDij0Uy9gTOKIEKV6OLUMB8uw0ezsD+eVPVtNc7f4GkeOUl+1lf4jQDKZMu9pXNkfP4p//ZmkKvwyR/WhiurSHLl3bVqP4sKuuHF3UIN9ffwxZ1YoL7QeO87+6xS696WnVqjOp9B2NvH75XeOrjC8600VbnPWo8T6nS8pXN8kNjf+zPEyza6pxuNqHH1WkuOq2mTX3Idmao2rnPqJvcqe62/9O8zFJmCNn2cmzbfzRmyaXbaD/0OiJ/ex5INC44s1SFIUHP/td5mpA/dWhaXx53KnTf67b2LTdiY/U/aeZFn6T8VH/Rq151PzB1918/Zl/2P/wR739ZNWQwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA38t/ZQErint9RokAAAAASUVORK5CYII=",
	}

	return &i
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

// GetRandomCreateCatalogImageContentOpts returns CreateImageCatalogContentOpts populated with random data
func GetRandomCreateCatalogImageContentOpts() *CreateImageCatalogContentOpts {
	c := CreateImageCatalogContentOpts{
		MediaID:   primitive.NewObjectID(),
		BrandID:   primitive.NewObjectID(),
		CatalogID: primitive.NewObjectID(),
		Label: &LabelOpts{
			Interests: faker.Lorem().Words(faker.RandomInt(1, 10)),
		},
	}
	return &c
}

// GetRandomCreateLiveStreamOpts returns CreateLiveStreamOpts with random data
func GetRandomCreateLiveStreamOpts() *CreateLiveStreamOpts {
	s := CreateLiveStreamOpts{
		Name: faker.Name().Name(),
		FeaturedImage: &Img{
			SRC: "https://deepakacademy.files.wordpress.com/2020/07/carryminatiyu.jpg",
		},
		StreamEndImage: &Img{
			SRC: "https://i2.wp.com/www.movieslantern.com/wp-content/uploads/2019/10/maxresdefault-170.jpg?fit=768%2C432&ssl=1",
		},
		InfluencerIDs: []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()},
		ScheduledAt:   time.Now().UTC().Add(7 * time.Minute),
	}
	for i := 0; i < faker.RandomInt(1, 10); i++ {
		s.CatalogIDs = append(s.CatalogIDs, primitive.NewObjectID())
	}
	return &s
}
