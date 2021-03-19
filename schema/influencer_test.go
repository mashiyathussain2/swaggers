package schema

import (
	"encoding/json"
	"go-app/server/validator"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_CreateInfluencerOpts(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    CreateInfluencerOpts
	}{
		{
			name: "[Ok]",
			json: string(`{
				"name": "test brand",
				"bio": "test bio",
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"external_links": [
					"https://youtube.com"
				],
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
			want: CreateInfluencerOpts{
				Name:          "test brand",
				Bio:           "test bio",
				ExternalLinks: []string{"https://youtube.com"},
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
			name: "[Ok] Without bio",
			json: string(`{
				"name": "test brand",
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"external_links": [
					"https://youtube.com"
				],
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
			want: CreateInfluencerOpts{
				Name:          "test brand",
				ExternalLinks: []string{"https://youtube.com"},
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
			name: "[Error] Without External Links",
			json: string(`{
				"name": "test brand",
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
			wantErr: true,
			err:     []string{"external_links is a required field"},
		},
		{
			name: "[Error] Empty External Links",
			json: string(`{
				"name": "test brand",
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"external_links": [
					
				],
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
			wantErr: true,
			err:     []string{"external_links must contain at least 1 item"},
		},
		{
			name: "[Ok] Without Instagram Social Media",
			json: string(`{
				"name": "test influencer",
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"external_links": [
					"https://youtube.com"
				],
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
			want: CreateInfluencerOpts{
				Name: "test influencer",
				CoverImg: &Img{
					SRC: "https://test.com/cover.png",
				},
				ExternalLinks: []string{"https://youtube.com"},
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
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"external_links": [
					"https://youtube.com"
				],
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
			want: CreateInfluencerOpts{
				Name: "test brand",
				CoverImg: &Img{
					SRC: "https://test.com/cover.png",
				},
				ExternalLinks: []string{"https://youtube.com"},
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
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"external_links": [
					"https://youtube.com"
				],
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc CreateInfluencerOpts
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

func Test_EditInfluencerOpts(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()
	id, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    EditInfluencerOpts
	}{
		{
			name: "[Ok] All fields",
			json: string(`{
				"id": "507f1f77bcf86cd799439011",
				"name": "test brand",
				"bio": "test bio",
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"external_links": [
					"https://youtube.com"
				],
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
			want: EditInfluencerOpts{
				ID:            id,
				Name:          "test brand",
				Bio:           "test bio",
				ExternalLinks: []string{"https://youtube.com"},
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
			name: "[Ok] Without bio",
			json: string(`{
				"id": "507f1f77bcf86cd799439011",
				"name": "test brand",
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"external_links": [
					"https://youtube.com"
				],
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
			want: EditInfluencerOpts{
				ID:            id,
				Name:          "test brand",
				ExternalLinks: []string{"https://youtube.com"},
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
			name: "[Ok] Without External Links",
			json: string(`{
				"id": "507f1f77bcf86cd799439011",
				"name": "test brand",
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
			want: EditInfluencerOpts{
				ID:   id,
				Name: "test brand",
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
			name: "[Error] Empty External Links",
			json: string(`{
				"id": "507f1f77bcf86cd799439011",
				"name": "test brand",
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"external_links": [
					
				],
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
			want: EditInfluencerOpts{
				ID:   id,
				Name: "test brand",
				CoverImg: &Img{
					SRC: "https://test.com/cover.png",
				},
				ExternalLinks: []string{},
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
				"name": "test influencer",
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"external_links": [
					"https://youtube.com"
				],
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
			want: EditInfluencerOpts{
				ID:   id,
				Name: "test influencer",
				CoverImg: &Img{
					SRC: "https://test.com/cover.png",
				},
				ExternalLinks: []string{"https://youtube.com"},
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
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"external_links": [
					"https://youtube.com"
				],
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
			want: EditInfluencerOpts{
				ID:   id,
				Name: "test brand",
				CoverImg: &Img{
					SRC: "https://test.com/cover.png",
				},
				ExternalLinks: []string{"https://youtube.com"},
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
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"external_links": [
					"https://youtube.com"
				],
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
			name: "[Error] Without ID",
			json: string(`{
				"name": "test brand",
				"cover_img": {
					"src": "https://test.com/cover.png"
				},
				"external_links": [
					"https://youtube.com"
				],
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
			wantErr: true,
			err:     []string{"id is a required field"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc EditInfluencerOpts
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
