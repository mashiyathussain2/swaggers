package schema

import (
	"encoding/json"
	"go-app/server/validator"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateBrandOpts(t *testing.T) {
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
				"description": "test deescription 1",
				"website_link": "https://link.testbrand.com",
				"fulfillment_email": "fulfill@testbrand.com"
			}`),
			wantErr: false,
			want: CreateBrandOpts{
				Name:             "test brand",
				RegisteredName:   "test brand pvt ltd",
				Description:      "test deescription 1",
				WebsiteLink:      "https://link.testbrand.com",
				FulfillmentEmail: "fulfill@testbrand.com",
			},
		},
		{
			name: "[ERROR] Without name",
			json: string(`{
				"registered_name": "test brand pvt ltd",
				"description": "test deescription 1",
				"website_link": "https://link.testbrand.com",
				"fulfillment_email": "fulfill@testbrand.com"
			}`),
			wantErr: true,
			err:     []string{"name is a required field"},
		},
		{
			name: "[ERROR] Without description",
			json: string(`{
				"name": "test brand",
				"registered_name": "test brand pvt ltd",
				"website_link": "https://link.testbrand.com",
				"fulfillment_email": "fulfill@testbrand.com"
			}`),
			wantErr: true,
			err:     []string{"description is a required field"},
		},
		{
			name: "[ERROR] Without invalid fulfillment email",
			json: string(`{
				"name": "test brand",
				"registered_name": "test brand pvt ltd",
				"description": "test deescription 1",
				"website_link": "https://link.testbrand.com",
				"fulfillment_email": "@testbrand.com"
			}`),
			wantErr: true,
			err:     []string{"fulfillment_email must be a valid email address"},
		},
		{
			name: "[ERROR] Without fulfillment email",
			json: string(`{
				"name": "test brand",
				"registered_name": "test brand pvt ltd",
				"description": "test deescription 1",
				"website_link": "https://link.testbrand.com"
			}`),
			wantErr: true,
			err:     []string{"fulfillment_email is a required field"},
		},
		{
			name: "[Ok] Without url",
			json: string(`{
				"name": "test brand",
				"description": "test deescription 1",
				"registered_name": "test brand pvt ltd",
				"fulfillment_email": "fulfill@testbrand.com"
			}`),
			wantErr: false,
			want: CreateBrandOpts{
				Name:             "test brand",
				RegisteredName:   "test brand pvt ltd",
				Description:      "test deescription 1",
				FulfillmentEmail: "fulfill@testbrand.com",
			},
		},
		{
			name: "[ERROR] With invalid url starting without http:// or https://",
			json: string(`{
				"name": "test brand",
				"description": "test deescription 1",
				"registered_name": "test brand pvt ltd",
				"fulfillment_email": "fulfill@testbrand.com",
				"website_link": "testbrand"
			}`),
			wantErr: true,
			err:     []string{"website_link must be a valid URL"},
		},
		{
			name: "[Ok] With valid url starting without http:// or https://",
			json: string(`{
				"name": "test brand",
				"description": "test deescription 1",
				"registered_name": "test brand pvt ltd",
				"fulfillment_email": "fulfill@testbrand.com",
				"website_link": "testbrand.com"
			}`),
			wantErr: false,
			want: CreateBrandOpts{
				Name:             "test brand",
				RegisteredName:   "test brand pvt ltd",
				Description:      "test deescription 1",
				FulfillmentEmail: "fulfill@testbrand.com",
				WebsiteLink:      "testbrand.com",
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
				assert.Equal(t, tt.want, sc)
			}
		})
	}
}

func TestEditBrandOpts(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()

	bID, _ := primitive.ObjectIDFromHex("6016a5baf83153e356354c34")
	bName := "test brand edited"
	bDesc := "test brand description edited"
	bLink := "https://link2.brand.com"
	bEmail := "email2.brand.com"

	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    EditBrandOpts
	}{
		{
			name: "[Ok] Edit Name",
			json: string(`{
				"id": "6016a5baf83153e356354c34",
				"name": "test brand edited"
			}`),
			wantErr: false,
			want: EditBrandOpts{
				ID:   bID,
				Name: bName,
			},
		},
		{
			name: "[Ok] Edit Description",
			json: string(`{
				"id": "6016a5baf83153e356354c34",
				"description": "test brand description edited"
			}`),
			wantErr: false,
			want: EditBrandOpts{
				ID:          bID,
				Description: bDesc,
			},
		},
		{
			name: "[Ok] Edit WebsiteLink",
			json: string(`{
				"id": "6016a5baf83153e356354c34",
				"website_link": "https://link2.brand.com"
			}`),
			wantErr: false,
			want: EditBrandOpts{
				ID:          bID,
				WebsiteLink: bLink,
			},
		},
		{
			name: "[ERROR] Edit WebsiteLink with invalid url",
			json: string(`{
				"id": "6016a5baf83153e356354c34",
				"website_link": "link2brandcom"
			}`),
			wantErr: true,
			err:     []string{"Key: 'EditBrandOpts.website_link' Error:Field validation for 'website_link' failed on the 'url|isdefault' tag"},
		},
		{
			name: "[Ok] Edit Fulfillment Email",
			json: string(`{
				"id": "6016a5baf83153e356354c34",
				"fulfillment_email": "email2.brand.com"
			}`),
			wantErr: false,
			want: EditBrandOpts{
				ID:               bID,
				FulfillmentEmail: bEmail,
			},
		},
		{
			name: "[Error] Edit Fulfillment Email with invalid email",
			json: string(`{
				"id": "6016a5baf83153e356354c34",
				"fulfillment_email": "brandcom.com"
			}`),
			wantErr: true,
			err:     []string{"Key: 'EditBrandOpts.fulfillment_email' Error:Field validation for 'fulfillment_email' failed on the 'email|isdefault' tag"},
		},
		{
			name: "[Ok] Multiple Fields (name, website_link, fulfillment_email)",
			json: string(`{
				"id": "6016a5baf83153e356354c34",
				"name": "test brand edited",
				"fulfillment_email": "email2.brand.com",
				"website_link": "https://link2.brand.com"
			}`),
			wantErr: false,
			want: EditBrandOpts{
				ID:               bID,
				FulfillmentEmail: bEmail,
				Name:             bName,
				WebsiteLink:      bLink,
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
				for i, e := range errs {
					assert.Equal(t, tt.err[i], e.Error())
				}
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want, sc)
			}
		})
	}
}
