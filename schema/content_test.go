package schema

import (
	"encoding/json"
	"go-app/server/validator"
	"testing"

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
				}
				}`),
			want: CreatePebbleOpts{
				Caption:       "Sample caption",
				InfluencerIDs: []primitive.ObjectID{id1, id2},
				BrandIDs:      []primitive.ObjectID{id3, id4},
				Label: &LabelOpts{
					Interests: []string{"fashion"},
					Gender:    []string{"M", "O"},
				},
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
				}
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
				}
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
				}
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
				}
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
				}
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
				}
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
