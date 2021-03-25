package schema

import (
	"encoding/json"
	"go-app/server/validator"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateCategoryOpts(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()

	pid, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2611")
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    CreateCategoryOpts
	}{
		{
			name: "[Ok] Without ParentID",
			json: string(`{
				"name": "Smartphones",
				"thumbnail": {
					"src": "https://images-eu.ssl-images-amazon.com/images/G/31/img18/Wireless/Catpage/BrandFarm/liwuwe_2018-05-07T11-25_f0461b_1113497_350x100_gps_cn_2.jpg"
				},
				"featured_image": {
					"src": "https://m.media-amazon.com/images/G/31/img20/Wireless/Apple/iPhone12/RiverImages/IN_r1307_r1306_Marketing_Page_L_FFH-1500_03._CB419228452_.jpg"
				},
				"is_main": true
			}`),
			wantErr: false,
			want: CreateCategoryOpts{
				Name: "Smartphones",
				Thumbnail: &Img{
					SRC: "https://images-eu.ssl-images-amazon.com/images/G/31/img18/Wireless/Catpage/BrandFarm/liwuwe_2018-05-07T11-25_f0461b_1113497_350x100_gps_cn_2.jpg",
				},
				FeaturedImage: &Img{
					SRC: "https://m.media-amazon.com/images/G/31/img20/Wireless/Apple/iPhone12/RiverImages/IN_r1307_r1306_Marketing_Page_L_FFH-1500_03._CB419228452_.jpg",
				},
				IsMain: true,
			},
		},
		{
			name: "[Ok] IsMain false",
			json: string(`{
				"name": "Smartphones",
				"thumbnail": {
					"src": "https://images-eu.ssl-images-amazon.com/images/G/31/img18/Wireless/Catpage/BrandFarm/liwuwe_2018-05-07T11-25_f0461b_1113497_350x100_gps_cn_2.jpg"
				},
				"featured_image": {
					"src": "https://m.media-amazon.com/images/G/31/img20/Wireless/Apple/iPhone12/RiverImages/IN_r1307_r1306_Marketing_Page_L_FFH-1500_03._CB419228452_.jpg"
				},
				"is_main": false
			}`),
			wantErr: false,
			want: CreateCategoryOpts{
				Name: "Smartphones",
				Thumbnail: &Img{
					SRC: "https://images-eu.ssl-images-amazon.com/images/G/31/img18/Wireless/Catpage/BrandFarm/liwuwe_2018-05-07T11-25_f0461b_1113497_350x100_gps_cn_2.jpg",
				},
				FeaturedImage: &Img{
					SRC: "https://m.media-amazon.com/images/G/31/img20/Wireless/Apple/iPhone12/RiverImages/IN_r1307_r1306_Marketing_Page_L_FFH-1500_03._CB419228452_.jpg",
				},
				IsMain: false,
			},
		},
		{
			name: "[OK] Without featured image",
			json: string(`{
				"name": "Smartphones",
				"thumbnail": {
					"src": "http://m.media-amazon.com/images/G/31/img20/Wireless/Apple/iPhone12/RiverImages/IN_r1307_r1306_Marketing_Page_L_FFH-1500_03._CB419228452_.jpg"
				},
				"parent_id": "5e8821fe1108c87837ef2611",
				"is_main": false
			}`),
			wantErr: false,
			want: CreateCategoryOpts{
				Name: "Smartphones",
				Thumbnail: &Img{
					SRC: "http://m.media-amazon.com/images/G/31/img20/Wireless/Apple/iPhone12/RiverImages/IN_r1307_r1306_Marketing_Page_L_FFH-1500_03._CB419228452_.jpg",
				},
				IsMain:   false,
				ParentID: pid,
			},
		},
		{
			name: "[OK] Without thumbnail",
			json: string(`{
				"name": "Smartphones",
				"featured_image": {
					"src": "http://m.media-amazon.com/images/G/31/img20/Wireless/Apple/iPhone12/RiverImages/IN_r1307_r1306_Marketing_Page_L_FFH-1500_03._CB419228452_.jpg"
				},
				"parent_id": "5e8821fe1108c87837ef2611",
				"is_main": false
			}`),
			wantErr: false,
			want: CreateCategoryOpts{
				Name: "Smartphones",
				FeaturedImage: &Img{
					SRC: "http://m.media-amazon.com/images/G/31/img20/Wireless/Apple/iPhone12/RiverImages/IN_r1307_r1306_Marketing_Page_L_FFH-1500_03._CB419228452_.jpg",
				},
				IsMain:   false,
				ParentID: pid,
			},
		},
		{
			name: "[Ok] With ParentID",
			json: string(`{
				"name": "Smartphones",
				"thumbnail": {
					"src": "https://images-eu.ssl-images-amazon.com/images/G/31/img18/Wireless/Catpage/BrandFarm/liwuwe_2018-05-07T11-25_f0461b_1113497_350x100_gps_cn_2.jpg"
				},
				"featured_image": {
					"src": "https://m.media-amazon.com/images/G/31/img20/Wireless/Apple/iPhone12/RiverImages/IN_r1307_r1306_Marketing_Page_L_FFH-1500_03._CB419228452_.jpg"
				},
				"parent_id": "5e8821fe1108c87837ef2611",
				"is_main": false
			}`),
			wantErr: false,
			want: CreateCategoryOpts{
				Name: "Smartphones",
				Thumbnail: &Img{
					SRC: "https://images-eu.ssl-images-amazon.com/images/G/31/img18/Wireless/Catpage/BrandFarm/liwuwe_2018-05-07T11-25_f0461b_1113497_350x100_gps_cn_2.jpg",
				},
				FeaturedImage: &Img{
					SRC: "https://m.media-amazon.com/images/G/31/img20/Wireless/Apple/iPhone12/RiverImages/IN_r1307_r1306_Marketing_Page_L_FFH-1500_03._CB419228452_.jpg",
				},
				IsMain:   false,
				ParentID: pid,
			},
		},
		{
			name: "[ERROR] With Invalid featured image SRC",
			json: string(`{
				"name": "Smartphones",
				"thumbnail": {
					"src": "https://images-eu.ssl-images-amazon.com/images/G/31/img18/Wireless/Catpage/BrandFarm/liwuwe_2018-05-07T11-25_f0461b_1113497_350x100_gps_cn_2.jpg"
				},
				"featured_image": {
					"src": "m.media-amazon.com/images/G/31/img20/Wireless/Apple/iPhone12/RiverImages/IN_r1307_r1306_Marketing_Page_L_FFH-1500_03._CB419228452_.jpg"
				},
				"parent_id": "5e8821fe1108c87837ef2611",
				"is_main": false
			}`),
			wantErr: true,
			err:     []string{"src must be a valid URL"},
		},
		{
			name: "[ERROR] With Invalid thumbnail SRC",
			json: string(`{
				"name": "Smartphones",
				"thumbnail": {
					"src": "images-eu.ssl-images-amazon.com/images/G/31/img18/Wireless/Catpage/BrandFarm/liwuwe_2018-05-07T11-25_f0461b_1113497_350x100_gps_cn_2.jpg"
				},
				"featured_image": {
					"src": "http://m.media-amazon.com/images/G/31/img20/Wireless/Apple/iPhone12/RiverImages/IN_r1307_r1306_Marketing_Page_L_FFH-1500_03._CB419228452_.jpg"
				},
				"parent_id": "5e8821fe1108c87837ef2611",
				"is_main": false
			}`),
			wantErr: true,
			err:     []string{"src must be a valid URL"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s CreateCategoryOpts
			err := json.Unmarshal([]byte(tt.json), &s)
			assert.Nil(t, err)
			errs := tv.Validate(&s)
			if tt.wantErr {
				assert.Len(t, errs, len(tt.err))
				assert.Equal(t, errs[0].Error(), tt.err[0])
			}
			if !tt.wantErr {
				assert.Len(t, errs, 0)
				assert.Equal(t, tt.want, s)
			}
		})
	}

}

func TestEditCategoryOpts(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()

	id, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2611")
	tbool := true
	fbool := false

	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    EditCategoryOpts
	}{
		{
			name: "[Ok] Name",
			json: string(`{
				"id": "5e8821fe1108c87837ef2611",
				"name": "Mobiles & Smartphones"
			}`),
			wantErr: false,
			want: EditCategoryOpts{
				ID:   id,
				Name: "Mobiles & Smartphones",
			},
		},
		{
			name: "[Ok] Thumbnail",
			json: string(`{
				"id": "5e8821fe1108c87837ef2611",
				"thumbnail": {
					"src": "https://thumbnail-edit.com/pic.png"
				}
			}`),
			wantErr: false,
			want: EditCategoryOpts{
				ID: id,
				Thumbnail: &Img{
					SRC: "https://thumbnail-edit.com/pic.png",
				},
			},
		},
		{
			name: "[Error] Invalid Thumbnail URL",
			json: string(`{
				"id": "5e8821fe1108c87837ef2611",
				"thumbnail": {
					"src": "thumbnail-edit.com/pic.png"
				}
			}`),
			wantErr: true,
			err:     []string{"src must be a valid URL"},
		},
		{
			name: "[Ok] Featured Image",
			json: string(`{
				"id": "5e8821fe1108c87837ef2611",
				"featured_image": {
					"src": "https://thumbnail-edit.com/pic.png"
				}
			}`),
			wantErr: false,
			want: EditCategoryOpts{
				ID: id,
				FeaturedImage: &Img{
					SRC: "https://thumbnail-edit.com/pic.png",
				},
			},
		},
		{
			name: "[Error] Invalid Featured Image URL",
			json: string(`{
				"id": "5e8821fe1108c87837ef2611",
				"featured_image": {
					"src": "thumbnail-edit.com"
				}
			}`),
			wantErr: true,
			err:     []string{"src must be a valid URL"},
		},
		{
			name: "[Ok] Is Main True",
			json: string(`{
				"id": "5e8821fe1108c87837ef2611",
				"is_main": true
			}`),
			wantErr: false,
			want: EditCategoryOpts{
				ID:     id,
				IsMain: &tbool,
			},
		},
		{
			name: "[Ok] Is Main False",
			json: string(`{
				"id": "5e8821fe1108c87837ef2611",
				"is_main": false
			}`),
			wantErr: false,
			want: EditCategoryOpts{
				ID:     id,
				IsMain: &fbool,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s EditCategoryOpts
			err := json.Unmarshal([]byte(tt.json), &s)
			assert.Nil(t, err)
			errs := tv.Validate(&s)
			if tt.wantErr {
				assert.Len(t, errs, len(tt.err))
				assert.Equal(t, errs[0].Error(), tt.err[0])
			}
			if !tt.wantErr {
				assert.Len(t, errs, 0)
				assert.Equal(t, tt.want, s)
			}
		})
	}

}
