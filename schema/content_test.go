package schema

import (
	"encoding/json"
	"go-app/model"
	"go-app/server/validator"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreatePebbleOpts(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()
	id1, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	id2, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439012")
	id3, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439013")
	id4, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439014")
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    CreatePebbleOpts
	}{
		{
			name: "[Ok]",
			json: string(`{
				"caption": "Sample caption",
				"influencer_ids": ["507f1f77bcf86cd799439011", "507f1f77bcf86cd799439012"],
				"brand_ids": ["507f1f77bcf86cd799439013", "507f1f77bcf86cd799439014"],
				"label": {
					"interests": ["fashion"],
					"gender" : ["M", "O"]
				},
				"file_name": "test.mp4"
			}`),
			want: CreatePebbleOpts{
				Caption:       "Sample caption",
				InfluencerIDs: []primitive.ObjectID{id1, id2},
				BrandIDs:      []primitive.ObjectID{id3, id4},
				Label: &LabelOpts{
					Interests: []string{"fashion"},
					Gender:    []string{"M", "O"},
				},
				FileName: "test.mp4",
			},
		}, {
			name: "[Ok] with Optional fields",
			json: string(`{
				"caption": "Sample caption",
				"influencer_ids": ["507f1f77bcf86cd799439011", "507f1f77bcf86cd799439012"],
				"brand_ids": ["507f1f77bcf86cd799439013", "507f1f77bcf86cd799439014"],
				"catalog_ids": ["507f1f77bcf86cd799439013", "507f1f77bcf86cd799439014"],
				"label": {
					"interests": ["fashion"],
					"age_group": ["18-22"],
					"gender" : ["M", "O"]
				},
				"file_name": "test.mp4"
			}`),
			want: CreatePebbleOpts{
				Caption:       "Sample caption",
				InfluencerIDs: []primitive.ObjectID{id1, id2},
				BrandIDs:      []primitive.ObjectID{id3, id4},
				CatalogIDs:    []primitive.ObjectID{id3, id4},
				Label: &LabelOpts{
					Interests: []string{"fashion"},
					AgeGroup:  []string{"18-22"},
					Gender:    []string{"M", "O"},
				},
				FileName: "test.mp4",
			},
		}, {
			name: "[Error] Missing InfluencerIDs",
			json: string(`{
				"caption": "Sample caption",
				"influencer_ids": [],
				"brand_ids": ["507f1f77bcf86cd799439013"],
				"label": {
					"interests": ["fashion"],
					"gender" : ["M", "O"]
				},
				"file_name": "test.mp4"
			}`),
			wantErr: true,
			err:     []string{"influencer_ids must contain at least 1 item"},
		}, {
			name: "[Error] Missing Brand IDs",
			json: string(`{
				"caption": "Sample caption",
				"influencer_ids": ["507f1f77bcf86cd799439013"],
				"brand_ids": [],
				"label": {
					"interests": ["fashion"],
					"gender" : ["M", "O"]
				},
				"file_name": "test.mp4"
			}`),
			wantErr: true,
			err:     []string{"brand_ids must contain at least 1 item"},
		}, {
			name: "[Error] Interest Label is missing",
			json: string(`{
				"caption": "Sample caption",
				"influencer_ids": ["507f1f77bcf86cd799439013"],
				"brand_ids": ["507f1f77bcf86cd799439014"],
				"label": {
					"interests": [],
					"gender" : ["M", "O"]
				},
				"file_name": "test.mp4"
			}`),
			wantErr: true,
			err:     []string{"interests must contain at least 1 item"},
		}, {
			name: "[Error] Gender label missing",
			json: string(`{
				"caption": "Sample caption",
				"influencer_ids": ["507f1f77bcf86cd799439013"],
				"brand_ids": ["507f1f77bcf86cd799439014"],
				"label": {
					"interests": ["Beauty"],
					"gender" : []
				},
				"file_name": "test.mp4"
			}`),
			wantErr: true,
			err:     []string{"gender must contain at least 1 item"},
		}, {
			name: "[Error] Unrecognized gender label",
			json: string(`{
				"caption": "Sample caption",
				"influencer_ids": ["507f1f77bcf86cd799439013"],
				"brand_ids": ["507f1f77bcf86cd799439014"],
				"label": {
					"interests": ["Beauty"],
					"gender" : ["Z"]
				},
				"file_name": "test.mp4"
			}`),
			wantErr: true,
			err:     []string{"gender[0] must be one of [M F O]"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc CreatePebbleOpts
			err := json.Unmarshal([]byte(tt.json), &sc)
			assert.Nil(t, err)
			errs := tv.Validate(&sc)
			if tt.wantErr {
				assert.Len(t, errs, len(tt.err))
				assert.Equal(t, errs[0].Error(), tt.err[0])
			}
			if !tt.wantErr {
				assert.Len(t, errs, 0)
				assert.Equal(t, tt.want, sc)
			}
		})
	}
}

func TestEditPebbleOpts(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()
	id1, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	id2, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439012")
	id3, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439013")
	id4, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439014")
	trueBool := true
	falseBool := false

	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    EditPebbleOpts
	}{
		{
			name: "[Ok] With all fields",
			json: string(`{
				"id": "507f1f77bcf86cd799439011",
				"caption": "Sample caption [edited]",
				"influencer_ids": ["507f1f77bcf86cd799439012"],
				"brand_ids": ["507f1f77bcf86cd799439013", "507f1f77bcf86cd799439014"],
				"catalog_ids": ["507f1f77bcf86cd799439013", "507f1f77bcf86cd799439014"],
				"label": {
					"interests": ["fashion", "men"],
					"gender" : ["M", "F"],
					"age_group": ["24-40", "40-60"]
				},
				"is_active": true
			}`),
			want: EditPebbleOpts{
				ID:            id1,
				Caption:       "Sample caption [edited]",
				InfluencerIDs: []primitive.ObjectID{id2},
				BrandIDs:      []primitive.ObjectID{id3, id4},
				CatalogIDs:    []primitive.ObjectID{id3, id4},
				Label: &EditLabelOpts{
					Interests: []string{"fashion", "men"},
					Gender:    []string{"M", "F"},
					AgeGroup:  []string{"24-40", "40-60"},
				},
				IsActive: &trueBool,
			},
		},
		{
			name: "[Ok] With IsActive False",
			json: string(`{
				"id": "507f1f77bcf86cd799439011",
				"caption": "Sample caption [edited]",
				"influencer_ids": ["507f1f77bcf86cd799439012"],
				"brand_ids": ["507f1f77bcf86cd799439013", "507f1f77bcf86cd799439014"],
				"catalog_ids": ["507f1f77bcf86cd799439013", "507f1f77bcf86cd799439014"],
				"label": {
					"interests": ["fashion", "men"],
					"gender" : ["M", "F"],
					"age_group": ["24-40", "40-60"]
				},
				"is_active": false
			}`),
			want: EditPebbleOpts{
				ID:            id1,
				Caption:       "Sample caption [edited]",
				InfluencerIDs: []primitive.ObjectID{id2},
				BrandIDs:      []primitive.ObjectID{id3, id4},
				CatalogIDs:    []primitive.ObjectID{id3, id4},
				Label: &EditLabelOpts{
					Interests: []string{"fashion", "men"},
					Gender:    []string{"M", "F"},
					AgeGroup:  []string{"24-40", "40-60"},
				},
				IsActive: &falseBool,
			},
		},
		{
			name: "[Ok] Without IsActive",
			json: string(`{
				"id": "507f1f77bcf86cd799439011",
				"caption": "Sample caption [edited]",
				"influencer_ids": ["507f1f77bcf86cd799439012"],
				"brand_ids": ["507f1f77bcf86cd799439013", "507f1f77bcf86cd799439014"],
				"catalog_ids": ["507f1f77bcf86cd799439013", "507f1f77bcf86cd799439014"],
				"label": {
					"interests": ["fashion", "men"],
					"gender" : ["M", "F"],
					"age_group": ["24-40", "40-60"]
				}
			}`),
			want: EditPebbleOpts{
				ID:            id1,
				Caption:       "Sample caption [edited]",
				InfluencerIDs: []primitive.ObjectID{id2},
				BrandIDs:      []primitive.ObjectID{id3, id4},
				CatalogIDs:    []primitive.ObjectID{id3, id4},
				Label: &EditLabelOpts{
					Interests: []string{"fashion", "men"},
					Gender:    []string{"M", "F"},
					AgeGroup:  []string{"24-40", "40-60"},
				},
				IsActive: nil,
			},
		},
		{
			name: "[Ok] With only Interests in Label",
			json: string(`{
				"id": "507f1f77bcf86cd799439011",
				"caption": "Sample caption [edited]",
				"influencer_ids": ["507f1f77bcf86cd799439012"],
				"brand_ids": ["507f1f77bcf86cd799439013", "507f1f77bcf86cd799439014"],
				"catalog_ids": ["507f1f77bcf86cd799439013", "507f1f77bcf86cd799439014"],
				"label": {
					"interests": ["fashion", "men"]
				}
			}`),
			want: EditPebbleOpts{
				ID:            id1,
				Caption:       "Sample caption [edited]",
				InfluencerIDs: []primitive.ObjectID{id2},
				BrandIDs:      []primitive.ObjectID{id3, id4},
				CatalogIDs:    []primitive.ObjectID{id3, id4},
				Label: &EditLabelOpts{
					Interests: []string{"fashion", "men"},
					// Gender:    []string{"M", "F"},
					// AgeGroup:  []string{"24-40", "40-60"},
				},
				IsActive: nil,
			},
		},
		{
			name: "[Ok] With only Interests in Label",
			json: string(`{
				"id": "507f1f77bcf86cd799439011",
				"caption": "Sample caption [edited]",
				"influencer_ids": ["507f1f77bcf86cd799439012"],
				"brand_ids": ["507f1f77bcf86cd799439013", "507f1f77bcf86cd799439014"],
				"catalog_ids": ["507f1f77bcf86cd799439013", "507f1f77bcf86cd799439014"],
				"label": {
					"gender" : ["M", "F"],
					"age_group": ["24-40", "40-60"]
				}
			}`),
			want: EditPebbleOpts{
				ID:            id1,
				Caption:       "Sample caption [edited]",
				InfluencerIDs: []primitive.ObjectID{id2},
				BrandIDs:      []primitive.ObjectID{id3, id4},
				CatalogIDs:    []primitive.ObjectID{id3, id4},
				Label: &EditLabelOpts{
					// Interests: []string{"fashion", "men"},
					Gender:   []string{"M", "F"},
					AgeGroup: []string{"24-40", "40-60"},
				},
				IsActive: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc EditPebbleOpts
			err := json.Unmarshal([]byte(tt.json), &sc)
			assert.Nil(t, err)
			errs := tv.Validate(&sc)
			if tt.wantErr {
				assert.Len(t, errs, len(tt.err))
				assert.Equal(t, errs[0].Error(), tt.err[0])
			}
			if !tt.wantErr {
				assert.Len(t, errs, 0)
				assert.Equal(t, tt.want, sc)
			}
		})
	}
}

func TestGetContentFilter(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()

	id1, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	id2, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439012")
	id3, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439013")
	trueBool := true
	falseBool := false
	str1 := "2020-02-08T15:04:05+00:00"
	from, _ := time.Parse(time.RFC3339, str1)
	str2 := "2020-02-16T15:04:05+00:00"
	to, _ := time.Parse(time.RFC3339, str2)
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    GetContentFilter
	}{
		{
			name: "[Ok]",
			json: string(`{
				"is_active": false,
				"is_processed": true,
				"media_type": "video",
				"type": "pebble",
				"brand_ids": [
					"507f1f77bcf86cd799439011"
				],
				"catalog_ids": [
					"507f1f77bcf86cd799439012",
					"507f1f77bcf86cd799439013"
				],
				"hashtags": [
					"#test",
					"#unitest"
				],
				"from": "2020-02-08T15:04:05+00:00",
				"to": "2020-02-16T15:04:05+00:00"
			}`),
			want: GetContentFilter{
				IsActive:    &falseBool,
				IsProcessed: &trueBool,
				MediaType:   model.VideoType,
				Type:        model.PebbleType,
				BrandIDs:    []primitive.ObjectID{id1},
				CatalogIDs:  []primitive.ObjectID{id2, id3},
				Hashtags:    []string{"#test", "#unitest"},
				From:        from,
				To:          to,
			},
		},
		{
			name: "[Error] Invalid media type",
			json: string(`{
				"is_active": false,
				"is_processed": true,
				"media_type": "audio",
				"type": "pebble",
				"brand_ids": [
					"507f1f77bcf86cd799439011"
				],
				"catalog_ids": [
					"507f1f77bcf86cd799439012",
					"507f1f77bcf86cd799439013"
				],
				"hashtags": [
					"#test",
					"#unitest"
				],
				"from": "2020-02-08T15:04:05+00:00",
				"to": "2020-02-16T15:04:05+00:00"
			}`),
			wantErr: true,
			err:     []string{"media_type must be one of [image video]"},
		},
		{
			name: "[Error] Invalid type",
			json: string(`{
				"is_active": false,
				"is_processed": true,
				"media_type": "video",
				"type": "user",
				"brand_ids": [
					"507f1f77bcf86cd799439011"
				],
				"catalog_ids": [
					"507f1f77bcf86cd799439012",
					"507f1f77bcf86cd799439013"
				],
				"hashtags": [
					"#test",
					"#unitest"
				],
				"from": "2020-02-08T15:04:05+00:00",
				"to": "2020-02-16T15:04:05+00:00"
			}`),
			wantErr: true,
			err:     []string{"type must be one of [pebble catalog_content]"},
		},
		{
			name: "[Ok] Without IsActive",
			json: string(`{
				"is_processed": true,
				"media_type": "image",
				"type": "catalog_content",
				"brand_ids": [
					"507f1f77bcf86cd799439011"
				],
				"catalog_ids": [
					"507f1f77bcf86cd799439012",
					"507f1f77bcf86cd799439013"
				],
				"hashtags": [
					"#test",
					"#unitest"
				],
				"from": "2020-02-08T15:04:05+00:00",
				"to": "2020-02-16T15:04:05+00:00"
			}`),
			want: GetContentFilter{
				IsProcessed: &trueBool,
				MediaType:   model.ImageType,
				Type:        model.CatalogContentType,
				BrandIDs:    []primitive.ObjectID{id1},
				CatalogIDs:  []primitive.ObjectID{id2, id3},
				Hashtags:    []string{"#test", "#unitest"},
				From:        from,
				To:          to,
			},
		},
		{
			name: "[Error] To is less than from",
			json: string(`{
				"is_processed": true,
				"media_type": "image",
				"type": "catalog_content",
				"brand_ids": [
					"507f1f77bcf86cd799439011"
				],
				"catalog_ids": [
					"507f1f77bcf86cd799439012",
					"507f1f77bcf86cd799439013"
				],
				"hashtags": [
					"#test",
					"#unitest"
				],
				"from": "2020-02-16T15:04:05+00:00",
				"to": "2020-02-08T15:04:05+00:00"
			}`),
			wantErr: true,
			err:     []string{"to must be greater than or equal to From"},
		},
		{
			name: "[Error] Without date filter",
			json: string(`{
				"is_processed": true,
				"media_type": "image",
				"type": "catalog_content",
				"brand_ids": [
					"507f1f77bcf86cd799439011"
				],
				"catalog_ids": [
					"507f1f77bcf86cd799439012",
					"507f1f77bcf86cd799439013"
				],
				"hashtags": [
					"#test",
					"#unitest"
				]
			}`),
			want: GetContentFilter{
				IsProcessed: &trueBool,
				MediaType:   model.ImageType,
				Type:        model.CatalogContentType,
				BrandIDs:    []primitive.ObjectID{id1},
				CatalogIDs:  []primitive.ObjectID{id2, id3},
				Hashtags:    []string{"#test", "#unitest"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc GetContentFilter
			err := json.Unmarshal([]byte(tt.json), &sc)
			assert.Nil(t, err)
			errs := tv.Validate(&sc)
			if tt.wantErr {
				assert.Len(t, errs, len(tt.err))
				assert.Equal(t, errs[0].Error(), tt.err[0])
			}
			if !tt.wantErr {
				assert.Len(t, errs, 0)
				assert.Equal(t, tt.want, sc)
			}
		})
	}
}
