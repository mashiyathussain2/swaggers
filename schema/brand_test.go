package schema

import (
	"encoding/json"
	"go-app/server/validator"
	"testing"

	"github.com/stretchr/testify/assert"
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
				assert.Len(t, errs, 0)
				assert.Equal(t, tt.want, sc)
			}
		})
	}
}
