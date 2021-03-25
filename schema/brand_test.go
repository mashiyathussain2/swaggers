package schema

import (
	"encoding/json"
	"go-app/server/validator"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_CreateBrandOpts(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    CreateBrandOpts
	}{
		{
			name: "[Ok]",
			json: string(`{
				"name": "test brand",
				"registered_name": "test brand pvt ltd",
				"fulfillment_email": "fulfullment@test.com",
				"fulfillment_cc_email": [
					"fulfullment1@test.com",
					"fulfullment2@test.com"
				],
				"domain": "test.com",
				"bio": "test bio",
				"website": "https://test.com",
				"logo": {
					"src": "https://test.com/logo.png"
				},
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"social_account": {
					"facebook": {
						"followers_count": 12000
					},
					"instagram": {
						"followers_count": 13000
					},
					"twitter": {
						"followers_count": 14000
					},
					"youtube": {
						"followers_count": 15000
					}
				}
			}`),
			want: CreateBrandOpts{
				Name:             "test brand",
				RegisteredName:   "test brand pvt ltd",
				FulfillmentEmail: "fulfullment@test.com",
				FulfillmentCCEmail: []string{
					"fulfullment1@test.com",
					"fulfullment2@test.com",
				},
				Domain:  "test.com",
				Website: "https://test.com",
				Bio:     "test bio",
				Logo: &Img{
					SRC: "https://test.com/logo.png",
				},
				CoverImg: &Img{
					SRC: "https://test.com/cover.png",
				},
				SocialAccount: &SocialAccountOpts{
					Facebook: &SocialMediaOpts{
						FollowersCount: 12000,
					},
					Instagram: &SocialMediaOpts{
						FollowersCount: 13000,
					},
					Twitter: &SocialMediaOpts{
						FollowersCount: 14000,
					},
					Youtube: &SocialMediaOpts{
						FollowersCount: 15000,
					},
				},
			},
		},
		{
			name: "[Ok] Without Instagram Social Media",
			json: string(`{
				"name": "test brand",
				"registered_name": "test brand pvt ltd",
				"fulfillment_email": "fulfullment@test.com",
				"fulfillment_cc_email": [
					"fulfullment1@test.com",
					"fulfullment2@test.com"
				],
				"domain": "test.com",
				"website": "https://test.com",
				"logo": {
					"src": "https://test.com/logo.png"
				},
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"social_account": {
					"facebook": {
						"followers_count": 12000
					},
					"twitter": {
						"followers_count": 14000
					},
					"youtube": {
						"followers_count": 15000
					}
				}
			}`),
			want: CreateBrandOpts{
				Name:             "test brand",
				RegisteredName:   "test brand pvt ltd",
				FulfillmentEmail: "fulfullment@test.com",
				FulfillmentCCEmail: []string{
					"fulfullment1@test.com",
					"fulfullment2@test.com",
				},
				Domain:  "test.com",
				Website: "https://test.com",
				Logo: &Img{
					SRC: "https://test.com/logo.png",
				},
				CoverImg: &Img{
					SRC: "https://test.com/cover.png",
				},
				SocialAccount: &SocialAccountOpts{
					Facebook: &SocialMediaOpts{
						FollowersCount: 12000,
					},
					Instagram: nil,
					Twitter: &SocialMediaOpts{
						FollowersCount: 14000,
					},
					Youtube: &SocialMediaOpts{
						FollowersCount: 15000,
					},
				},
			},
		},
		{
			name: "[Ok] With Instagram Follower Count is 0",
			json: string(`{
				"name": "test brand",
				"registered_name": "test brand pvt ltd",
				"fulfillment_email": "fulfullment@test.com",
				"fulfillment_cc_email": [
					"fulfullment1@test.com",
					"fulfullment2@test.com"
				],
				"domain": "test.com",
				"website": "https://test.com",
				"logo": {
					"src": "https://test.com/logo.png"
				},
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"social_account": {
					"facebook": {
						"followers_count": 12000
					},
					"instagram": {
						"followers_count": 0
					},
					"twitter": {
						"followers_count": 14000
					},
					"youtube": {
						"followers_count": 15000
					}
				}
			}`),
			want: CreateBrandOpts{
				Name:             "test brand",
				RegisteredName:   "test brand pvt ltd",
				FulfillmentEmail: "fulfullment@test.com",
				FulfillmentCCEmail: []string{
					"fulfullment1@test.com",
					"fulfullment2@test.com",
				},
				Domain:  "test.com",
				Website: "https://test.com",
				Logo: &Img{
					SRC: "https://test.com/logo.png",
				},
				CoverImg: &Img{
					SRC: "https://test.com/cover.png",
				},
				SocialAccount: &SocialAccountOpts{
					Facebook: &SocialMediaOpts{
						FollowersCount: 12000,
					},
					Instagram: &SocialMediaOpts{
						FollowersCount: 0,
					},
					Twitter: &SocialMediaOpts{
						FollowersCount: 14000,
					},
					Youtube: &SocialMediaOpts{
						FollowersCount: 15000,
					},
				},
			},
		},
		{
			name: "[Error] With Instagram Follower Count is less than 0",
			json: string(`{
				"name": "test brand",
				"registered_name": "test brand pvt ltd",
				"fulfillment_email": "fulfullment@test.com",
				"fulfillment_cc_email": [
					"fulfullment1@test.com",
					"fulfullment2@test.com"
				],
				"domain": "test.com",
				"website": "https://test.com",
				"logo": {
					"src": "https://test.com/logo.png"
				},
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"social_account": {
					"facebook": {
						"followers_count": 12000
					},
					"instagram": {
						"followers_count": -1
					},
					"twitter": {
						"followers_count": 14000
					},
					"youtube": {
						"followers_count": 15000
					}
				}
			}`),
			wantErr: true,
			err:     []string{"followers_count must be 0 or greater"},
		},
		{
			name: "[Error] With invalid website url",
			json: string(`{
				"name": "test brand",
				"registered_name": "test brand pvt ltd",
				"fulfillment_email": "fulfullment@test.com",
				"fulfillment_cc_email": [
					"fulfullment1@test.com",
					"fulfullment2@test.com"
				],
				"domain": "test.com",
				"website": "test.com",
				"logo": {
					"src": "https://test.com/logo.png"
				},
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"social_account": {
					"facebook": {
						"followers_count": 12000
					},
					"instagram": {
						"followers_count": 1200
					},
					"twitter": {
						"followers_count": 14000
					},
					"youtube": {
						"followers_count": 15000
					}
				}
			}`),
			wantErr: true,
			err:     []string{"website must be a valid URL"},
		},
		{
			name: "[Error] With invalid fulfillment_cc_email",
			json: string(`{
				"name": "test brand",
				"registered_name": "test brand pvt ltd",
				"fulfillment_email": "fulfullment@test.com",
				"fulfillment_cc_email": [
					"fulfullment1@test.com",
					"fulfullment"
				],
				"domain": "test.com",
				"website": "https://www.test.com",
				"logo": {
					"src": "https://test.com/logo.png"
				},
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"social_account": {
					"facebook": {
						"followers_count": 12000
					},
					"instagram": {
						"followers_count": 1200
					},
					"twitter": {
						"followers_count": 14000
					},
					"youtube": {
						"followers_count": 15000
					}
				}
			}`),
			wantErr: true,
			err:     []string{"fulfillment_cc_email[1] must be a valid email address"},
		},
		{
			name: "[Error] With invalid fulfillment_cc_email",
			json: string(`{
				"name": "test brand",
				"registered_name": "test brand pvt ltd",
				"fulfillment_email": "fulfullment_test.com",
				"fulfillment_cc_email": [
					"fulfullment1@test.com"
				],
				"domain": "test.com",
				"website": "https://www.test.com",
				"logo": {
					"src": "https://test.com/logo.png"
				},
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"social_account": {
					"facebook": {
						"followers_count": 12000
					},
					"instagram": {
						"followers_count": 1200
					},
					"twitter": {
						"followers_count": 14000
					},
					"youtube": {
						"followers_count": 15000
					}
				}
			}`),
			wantErr: true,
			err:     []string{"fulfillment_email must be a valid email address"},
		},
		{
			name: "[Ok] Without FulfillmentCCEmail",
			json: string(`{
				"name": "test brand",
				"registered_name": "test brand pvt ltd",
				"fulfillment_email": "fulfullment@test.com",
				"domain": "test.com",
				"website": "https://test.com",
				"logo": {
					"src": "https://test.com/logo.png"
				},
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"social_account": {
					"facebook": {
						"followers_count": 12000
					},
					"instagram": {
						"followers_count": 13000
					},
					"twitter": {
						"followers_count": 14000
					},
					"youtube": {
						"followers_count": 15000
					}
				}
			}`),
			want: CreateBrandOpts{
				Name:               "test brand",
				RegisteredName:     "test brand pvt ltd",
				FulfillmentEmail:   "fulfullment@test.com",
				FulfillmentCCEmail: nil,
				Domain:             "test.com",
				Website:            "https://test.com",
				Logo: &Img{
					SRC: "https://test.com/logo.png",
				},
				CoverImg: &Img{
					SRC: "https://test.com/cover.png",
				},
				SocialAccount: &SocialAccountOpts{
					Facebook: &SocialMediaOpts{
						FollowersCount: 12000,
					},
					Instagram: &SocialMediaOpts{
						FollowersCount: 13000,
					},
					Twitter: &SocialMediaOpts{
						FollowersCount: 14000,
					},
					Youtube: &SocialMediaOpts{
						FollowersCount: 15000,
					},
				},
			},
		},
		{
			name: "[Ok] Empty FulfillmentCCEmail",
			json: string(`{
				"name": "test brand",
				"registered_name": "test brand pvt ltd",
				"fulfillment_email": "fulfullment@test.com",
				"domain": "test.com",
				"website": "https://test.com",
				"fulfillment_cc_email": [],
				"logo": {
					"src": "https://test.com/logo.png"
				},
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"social_account": {
					"facebook": {
						"followers_count": 12000
					},
					"instagram": {
						"followers_count": 13000
					},
					"twitter": {
						"followers_count": 14000
					},
					"youtube": {
						"followers_count": 15000
					}
				}
			}`),
			want: CreateBrandOpts{
				Name:               "test brand",
				RegisteredName:     "test brand pvt ltd",
				FulfillmentEmail:   "fulfullment@test.com",
				FulfillmentCCEmail: []string{},
				Domain:             "test.com",
				Website:            "https://test.com",
				Logo: &Img{
					SRC: "https://test.com/logo.png",
				},
				CoverImg: &Img{
					SRC: "https://test.com/cover.png",
				},
				SocialAccount: &SocialAccountOpts{
					Facebook: &SocialMediaOpts{
						FollowersCount: 12000,
					},
					Instagram: &SocialMediaOpts{
						FollowersCount: 13000,
					},
					Twitter: &SocialMediaOpts{
						FollowersCount: 14000,
					},
					Youtube: &SocialMediaOpts{
						FollowersCount: 15000,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc CreateBrandOpts
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

func Test_EditBrandOpts(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()

	id, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    EditBrandOpts
	}{
		{
			name: "[Ok] All fields",
			json: string(`{
				"id": "507f1f77bcf86cd799439011",
				"name": "test brand",
				"registered_name": "test brand pvt ltd",
				"fulfillment_email": "fulfullment@test.com",
				"fulfillment_cc_email": [
					"fulfullment1@test.com",
					"fulfullment2@test.com"
				],

				"bio": "test bio",
				"domain": "test.com",
				"website": "https://test.com",
				"logo": {
					"src": "https://test.com/logo.png"
				},
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"social_account": {
					"facebook": {
						"followers_count": 12000
					},
					"instagram": {
						"followers_count": 13000
					},
					"twitter": {
						"followers_count": 14000
					},
					"youtube": {
						"followers_count": 15000
					}
				}
			}`),
			want: EditBrandOpts{
				ID:               id,
				Name:             "test brand",
				RegisteredName:   "test brand pvt ltd",
				FulfillmentEmail: "fulfullment@test.com",
				FulfillmentCCEmail: []string{
					"fulfullment1@test.com",
					"fulfullment2@test.com",
				},
				Bio:     "test bio",
				Domain:  "test.com",
				Website: "https://test.com",
				Logo: &Img{
					SRC: "https://test.com/logo.png",
				},
				CoverImg: &Img{
					SRC: "https://test.com/cover.png",
				},
				SocialAccount: &SocialAccountOpts{
					Facebook: &SocialMediaOpts{
						FollowersCount: 12000,
					},
					Instagram: &SocialMediaOpts{
						FollowersCount: 13000,
					},
					Twitter: &SocialMediaOpts{
						FollowersCount: 14000,
					},
					Youtube: &SocialMediaOpts{
						FollowersCount: 15000,
					},
				},
			},
		},
		{
			name: "[Ok] Without Instagram Social Media",
			json: string(`{
				"id": "507f1f77bcf86cd799439011",
				"name": "test brand",
				"registered_name": "test brand pvt ltd",
				"fulfillment_email": "fulfullment@test.com",
				"fulfillment_cc_email": [
					"fulfullment1@test.com",
					"fulfullment2@test.com"
				],
				"domain": "test.com",
				"website": "https://test.com",
				"logo": {
					"src": "https://test.com/logo.png"
				},
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"social_account": {
					"facebook": {
						"followers_count": 12000
					},
					"twitter": {
						"followers_count": 14000
					},
					"youtube": {
						"followers_count": 15000
					}
				}
			}`),
			want: EditBrandOpts{
				ID:               id,
				Name:             "test brand",
				RegisteredName:   "test brand pvt ltd",
				FulfillmentEmail: "fulfullment@test.com",
				FulfillmentCCEmail: []string{
					"fulfullment1@test.com",
					"fulfullment2@test.com",
				},
				Domain:  "test.com",
				Website: "https://test.com",
				Logo: &Img{
					SRC: "https://test.com/logo.png",
				},
				CoverImg: &Img{
					SRC: "https://test.com/cover.png",
				},
				SocialAccount: &SocialAccountOpts{
					Facebook: &SocialMediaOpts{
						FollowersCount: 12000,
					},
					Instagram: nil,
					Twitter: &SocialMediaOpts{
						FollowersCount: 14000,
					},
					Youtube: &SocialMediaOpts{
						FollowersCount: 15000,
					},
				},
			},
		},
		{
			name: "[Ok] With Instagram Follower Count is 0",
			json: string(`{
				"id": "507f1f77bcf86cd799439011",
				"name": "test brand",
				"registered_name": "test brand pvt ltd",
				"fulfillment_email": "fulfullment@test.com",
				"fulfillment_cc_email": [
					"fulfullment1@test.com",
					"fulfullment2@test.com"
				],
				"domain": "test.com",
				"website": "https://test.com",
				"logo": {
					"src": "https://test.com/logo.png"
				},
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"social_account": {
					"facebook": {
						"followers_count": 12000
					},
					"instagram": {
						"followers_count": 0
					},
					"twitter": {
						"followers_count": 14000
					},
					"youtube": {
						"followers_count": 15000
					}
				}
			}`),
			want: EditBrandOpts{
				ID:               id,
				Name:             "test brand",
				RegisteredName:   "test brand pvt ltd",
				FulfillmentEmail: "fulfullment@test.com",
				FulfillmentCCEmail: []string{
					"fulfullment1@test.com",
					"fulfullment2@test.com",
				},
				Domain:  "test.com",
				Website: "https://test.com",
				Logo: &Img{
					SRC: "https://test.com/logo.png",
				},
				CoverImg: &Img{
					SRC: "https://test.com/cover.png",
				},
				SocialAccount: &SocialAccountOpts{
					Facebook: &SocialMediaOpts{
						FollowersCount: 12000,
					},
					Instagram: &SocialMediaOpts{
						FollowersCount: 0,
					},
					Twitter: &SocialMediaOpts{
						FollowersCount: 14000,
					},
					Youtube: &SocialMediaOpts{
						FollowersCount: 15000,
					},
				},
			},
		},
		{
			name: "[Error] With Instagram Follower Count is less than 0",
			json: string(`{
				"id": "507f1f77bcf86cd799439011",
				"name": "test brand",
				"registered_name": "test brand pvt ltd",
				"fulfillment_email": "fulfullment@test.com",
				"fulfillment_cc_email": [
					"fulfullment1@test.com",
					"fulfullment2@test.com"
				],
				"domain": "test.com",
				"website": "https://test.com",
				"logo": {
					"src": "https://test.com/logo.png"
				},
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"social_account": {
					"facebook": {
						"followers_count": 12000
					},
					"instagram": {
						"followers_count": -1
					},
					"twitter": {
						"followers_count": 14000
					},
					"youtube": {
						"followers_count": 15000
					}
				}
			}`),
			wantErr: true,
			err:     []string{"followers_count must be 0 or greater"},
		},
		{
			name: "[Error] With invalid website url",
			json: string(`{
				"id": "507f1f77bcf86cd799439011",
				"website": "test.com"
			}`),
			wantErr: true,
			err:     []string{"Key: 'EditBrandOpts.website' Error:Field validation for 'website' failed on the 'isdefault|url' tag"},
		},
		{
			name: "[Error] With invalid fulfillment_cc_email",
			json: string(`{
				"id": "507f1f77bcf86cd799439011",
				"fulfillment_cc_email": [
					"fulfullment1@test.com",
					"fulfullment"
				]
			}`),
			wantErr: true,
			err:     []string{"fulfillment_cc_email[1] must be a valid email address"},
		},
		{
			name: "[Ok] Empty FulfillmentCCEmail",
			json: string(`{
				"id": "507f1f77bcf86cd799439011",
				"fulfillment_cc_email": []
			}`),
			want: EditBrandOpts{
				ID:                 id,
				FulfillmentCCEmail: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc EditBrandOpts
			err := json.Unmarshal([]byte(tt.json), &sc)
			assert.Nil(t, err)
			errs := tv.Validate(&sc)
			if tt.wantErr {
				assert.Len(t, errs, len(tt.err))
				assert.Equal(t, tt.err[0], errs[0].Error())
			}
			if !tt.wantErr {
				assert.Len(t, errs, 0)
				assert.Equal(t, tt.want, sc)
			}
		})
	}
}
