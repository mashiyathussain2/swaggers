package schema

import (
	"encoding/json"
	"go-app/model"
	"go-app/server/validator"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateCatalogOpts(t *testing.T) {
	t.Parallel()
	bID, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2611")
	cID, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2612")
	tv := validator.NewValidation()
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    CreateCatalogOpts
	}{
		{
			name: "[Ok] Without VariantType and Variant",
			json: string(`{
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"brand_id": "5e8821fe1108c87837ef2611",
				"description": "test description 1",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"base_price": 1299,
				"retail_price": 1099,
				"featured_image":{
					"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
				}
			}`),
			wantErr: false,
			want: CreateCatalogOpts{
				Name:        "test",
				CategoryID:  []primitive.ObjectID{cID},
				Description: "test description 1",
				BrandID:     bID,
				Keywords:    []string{"k1", "k2"},
				HSNCode:     "hsnCode1",
				BasePrice:   1299,
				RetailPrice: 1099,
				FeaturedImage: &Img{
					SRC: "https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31",
				},
			},
		},
		{
			name: "[Ok] With Filter attribute",
			json: string(`{
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"brand_id": "5e8821fe1108c87837ef2611",
				"description": "test description 1",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"base_price": 1299,
				"retail_price": 1099,
				"featured_image":{
					"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
				},
				"filter_attr": [{
					"name": "Color",
					"value": "Red"
				}]
			}`),
			wantErr: false,
			want: CreateCatalogOpts{
				Name:        "test",
				CategoryID:  []primitive.ObjectID{cID},
				Description: "test description 1",
				BrandID:     bID,
				Keywords:    []string{"k1", "k2"},
				HSNCode:     "hsnCode1",
				BasePrice:   1299,
				RetailPrice: 1099,
				FilterAttribute: []FilterAttribute{
					{
						Name:  "Color",
						Value: "Red",
					},
				},
				FeaturedImage: &Img{
					SRC: "https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31",
				},
			},
		},
		{
			name: "[Ok] With VariantType and Variant",
			json: string(`{
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"keywords":  ["k1", "k2"],
				"brand_id": "5e8821fe1108c87837ef2611",
				"hsn_code": "hsnCode1",
				"variant_type": "size",
				"base_price": 1299,
				"retail_price": 1099,
				"featured_image":{
					"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
				},
				"variants": [
					{
						"sku": "sku1",
						"attribute": "red",
						"unit":2
					}
				]
			}`),
			wantErr: false,
			want: CreateCatalogOpts{
				Name:        "test",
				CategoryID:  []primitive.ObjectID{cID},
				BrandID:     bID,
				Description: "test description 1",
				Keywords:    []string{"k1", "k2"},
				HSNCode:     "hsnCode1",
				VariantType: model.SizeType,
				BasePrice:   1299,
				RetailPrice: 1099,
				Variants: []CreateVariantOpts{
					{
						SKU:       "sku1",
						Attribute: "red",
						Unit:      2,
					},
				},
				FeaturedImage: &Img{
					SRC: "https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31",
				},
			},
		},
		{
			name: "[Error] Passing duplicate keywords",
			json: string(`{
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"brand_id": "5e8821fe1108c87837ef2611",
				"keywords":  ["k1", "k1"],
				"hsn_code": "hsnCode1",
				"base_price": 1299,
				"featured_image":{
					"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
				},
				"retail_price": 1099
			}`),
			wantErr: true,
			err:     []string{"keywords must contain unique values"},
		},
		{
			name: "[Error] Passing variant type only",
			json: string(`{
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"brand_id": "5e8821fe1108c87837ef2611",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"variant_type": "size",
				"base_price": 1299,
				"featured_image":{
					"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
				},
				"retail_price": 1099
			}`),
			wantErr: true,
			err: []string{
				"Key: 'CreateCatalogOpts.variant_type' Error:Field validation for 'variant_type' failed on the 'required_with_field' tag",
				// "Key: 'CreateCatalogOpts.variants' Error:Field validation for 'variants' failed on the 'required_with_field' tag",
			},
		},
		{
			name: "[Error] Passing variants only",
			json: string(`{
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"brand_id": "5e8821fe1108c87837ef2611",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"featured_image":{
					"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
				},
				"variants": [
					{
						"sku": "sku1",
						"attribute": "red",
						"unit":2
					}
				],
				"base_price": 1299,
				"retail_price": 1099
			}`),
			wantErr: true,
			err: []string{
				// "Key: 'CreateCatalogOpts.variant_type' Error:Field validation for 'variant_type' failed on the 'required_with_field' tag",
				"Key: 'CreateCatalogOpts.variant_type' Error:Field validation for 'variant_type' failed on the 'required_with_field' tag",
			},
		},
		{
			name: "[Error] Without Brand id",
			json: string(`{
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"featured_image":{
					"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
				},
				"base_price": 1299,
				"retail_price": 1099
			}`),
			wantErr: true,
			err: []string{
				// "Key: 'CreateCatalogOpts.variant_type' Error:Field validation for 'variant_type' failed on the 'required_with_field' tag",
				"brand_id is a required field",
			},
		},
		{
			name: "[Ok] With ETA",
			json: string(`{
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"brand_id": "5e8821fe1108c87837ef2611",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"featured_image":{
					"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
				},
				"eta": {
					"min": 1,
					"max": 7,
					"unit": "day"
				},
				"base_price": 1299,
				"retail_price": 1099
			}`),
			wantErr: false,
			want: CreateCatalogOpts{
				Name:        "test",
				CategoryID:  []primitive.ObjectID{cID},
				BrandID:     bID,
				Description: "test description 1",
				Keywords:    []string{"k1", "k2"},
				HSNCode:     "hsnCode1",
				ETA: &etaOpts{
					Min:  1,
					Max:  7,
					Unit: "day",
				},
				FeaturedImage: &Img{
					SRC: "https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31",
				},
				BasePrice:   1299,
				RetailPrice: 1099,
				// VariantType: model.SizeType,
				// Variants: []CreateVariantOpts{
				// 	{
				// 		SKU:         "sku1",
				// 		BasePrice:   1299,
				// 		RetailPrice: 1099,
				// 	},
				// },
			},
		},
		{
			name: "[Error] With invalid eta unit",
			json: string(`{
					"name": "test",
					"category_id": ["5e8821fe1108c87837ef2612"],
					"description": "test description 1",
					"brand_id": "5e8821fe1108c87837ef2611",
					"keywords":  ["k1", "k2"],
					"hsn_code": "hsnCode1",
					"eta": {
						"min": 1,
						"max": 7,
						"unit": "year"
					},
					"base_price": 1299,
					"retail_price": 1099,
					"featured_image":{
						"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
					}
				}`),
			wantErr: true,
			err:     []string{"unit must be one of [hour day month]"},
		},
		{
			name: "[Ok] With Specifications",
			json: string(`{
					"name": "test",
					"category_id": ["5e8821fe1108c87837ef2612"],
					"description": "test description 1",
					"brand_id": "5e8821fe1108c87837ef2611",
					"keywords":  ["k1", "k2"],
					"hsn_code": "hsnCode1",
					"specifications": [{
						"Name": "k1",
						"Value": "v1"
					},{
						"Name": "k2",
						"Value": "v2"
					}],
					"base_price": 1299,
					"retail_price": 1099,
					"featured_image":{
						"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
					}
				}`),
			wantErr: false,
			want: CreateCatalogOpts{
				Name:        "test",
				CategoryID:  []primitive.ObjectID{cID},
				BrandID:     bID,
				Description: "test description 1",
				Keywords:    []string{"k1", "k2"},
				HSNCode:     "hsnCode1",
				BasePrice:   1299,
				RetailPrice: 1099,
				FeaturedImage: &Img{
					SRC: "https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31",
				},
				Specifications: []specsOpts{
					{
						Name:  "k1",
						Value: "v1",
					},
					{
						Name:  "k2",
						Value: "v2",
					},
				},
			},
		},
		{
			name: "[Error] With empty Name[1] field specification",
			json: string(`{
					"name": "test",
					"category_id": ["5e8821fe1108c87837ef2612"],
					"description": "test description 1",
					"brand_id": "5e8821fe1108c87837ef2611",
					"keywords":  ["k1", "k2"],
					"hsn_code": "hsnCode1",
					"specifications": [{
						"Name": "",
						"Value": "v2"
					}],
					"featured_image":{
						"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
					},
					"base_price": 1299,
					"retail_price": 1099
				}`),
			wantErr: true,
			err:     []string{"name is a required field"},
		},
		{
			name: "[Error] With empty Value[0] field specification",
			json: string(`{
					"name": "test",
					"category_id": ["5e8821fe1108c87837ef2612"],
					"description": "test description 1",
					"brand_id": "5e8821fe1108c87837ef2611",
					"keywords":  ["k1", "k2"],
					"hsn_code": "hsnCode1",
					"specifications": [{
						"Name": "k1",
						"Value": ""
					},{
						"Name": "k2",
						"Value": "v2"
					}],
					"base_price": 1299,
					"retail_price": 1099,
					"featured_image":{
						"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
					}
				}`),
			wantErr: true,
			err:     []string{"value is a required field"},
		},
		{
			name: "[Error] Without category id",
			json: string(`{
					"name": "test",
					"description": "test description 1",
					"brand_id": "5e8821fe1108c87837ef2611",
					"keywords":  ["k1", "k2"],
					"hsn_code": "hsnCode1",
					"base_price": 1299,
					"retail_price": 1099,
					"featured_image":{
						"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
					}
				}`),
			wantErr: true,
			err:     []string{"category_id is a required field"},
		},
		{
			name: "[Error] Without base price",
			json: string(`{
					"name": "test",
					"description": "test description 1",
					"category_id": ["5e8821fe1108c87837ef2612"],
					"brand_id": "5e8821fe1108c87837ef2611",
					"keywords":  ["k1", "k2"],
					"hsn_code": "hsnCode1",
					"retail_price": 1099,
					"featured_image":{
						"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
					}
				}`),
			wantErr: true,
			err:     []string{"base_price must be greater than 0"},
		},
		{
			name: "[Error] With base price less than retail price",
			json: string(`{
					"name": "test",
					"description": "test description 1",
					"category_id": ["5e8821fe1108c87837ef2612"],
					"brand_id": "5e8821fe1108c87837ef2611",
					"keywords":  ["k1", "k2"],
					"hsn_code": "hsnCode1",
					"base_price": 999,
					"retail_price": 1099,
					"featured_image":{
						"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
					}
				}`),
			wantErr: true,
			err:     []string{"base_price must be greater than or equal to RetailPrice"},
		},
		{
			name: "[Error] Without retail price",
			json: string(`{
					"name": "test",
					"description": "test description 1",
					"category_id": ["5e8821fe1108c87837ef2612"],
					"brand_id": "5e8821fe1108c87837ef2611",
					"keywords":  ["k1", "k2"],
					"hsn_code": "hsnCode1",
					"base_price": 1099,
					"featured_image":{
						"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
					}
				}`),
			wantErr: true,
			err:     []string{"retail_price must be greater than 0"},
		},
		{
			name: "[Error] With retail price 0",
			json: string(`{
					"name": "test",
					"description": "test description 1",
					"category_id": ["5e8821fe1108c87837ef2612"],
					"brand_id": "5e8821fe1108c87837ef2611",
					"keywords":  ["k1", "k2"],
					"hsn_code": "hsnCode1",
					"base_price": 1299,
					"retail_price": 0,
					"featured_image":{
						"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
					}
				}`),
			wantErr: true,
			err:     []string{"retail_price must be greater than 0"},
		},
		{
			name: "[Error] Withempty array of category id",
			json: string(`{
					"name": "test",
					"category_id": [],
					"description": "test description 1",
					"brand_id": "5e8821fe1108c87837ef2611",
					"keywords":  ["k1", "k2"],
					"hsn_code": "hsnCode1",
					"base_price": 1299,
					"retail_price": 1099,
					"featured_image":{
						"src":"https://img-prod-cms-rt-microsoft-com.akamaized.net/cms/api/am/imageFileData/RE1Mu3b?ver=5c31"
					}
				}`),
			wantErr: true,
			err:     []string{"category_id must contain more than 0 items"},
		},
		{
			name: "[Error] Without featured image",
			json: string(`{
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"brand_id": "5e8821fe1108c87837ef2611",
				"description": "test description 1",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"base_price": 1299,
				"retail_price": 1099
				
			}`),
			wantErr: true,
			err:     []string{"featured_image is a required field"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc CreateCatalogOpts
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

func TestCreateVariantOpts(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    CreateVariantOpts
	}{
		{
			name: "[Ok]",
			json: string(`{
				"sku": "sku1",
				"attribute": "red"
			}`),
			wantErr: false,
			want: CreateVariantOpts{
				SKU:       "sku1",
				Attribute: "red",
			},
		},
		{
			name: "[Ok] With Attribute",
			json: string(`{
				"sku": "sku1",
				"attribute": "Red"
			}`),
			wantErr: false,
			want: CreateVariantOpts{
				SKU:       "sku1",
				Attribute: "Red",
			},
		},
		{
			name: "[Error] Empty SKU",
			json: string(`{
				"sku": "",
				"attribute": "red"
			}`),
			wantErr: true,
			err:     []string{"sku is a required field"},
		},
		{
			name: "[Error] No SKU",
			json: string(`{
				"attribute": "red"
			}`),
			wantErr: true,
			err:     []string{"sku is a required field"},
		},
		{
			name: "[Error] No Attribute",
			json: string(`{
				"sku": "red-1"
			}`),
			wantErr: true,
			err:     []string{"attribute is a required field"},
		},
		{
			name: "[Error] Empty Attribute",
			json: string(`{
				"sku": "red-1",
				"attribute": ""
			}`),
			wantErr: true,
			err:     []string{"attribute is a required field"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc CreateVariantOpts
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

func TestAddVariantOpts(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    AddVariantOpts
	}{
		{
			name: "[Ok]",
			json: string(`{
				"variant_type": "size",
				"sku": "sku1",
				"attribute": "red"
			}`),
			wantErr: false,
			want: AddVariantOpts{
				VariantType: "size",
				SKU:         "sku1",
				Attribute:   "red",
			},
		},
		{
			name: "[Ok] With Attribute",
			json: string(`{
				"sku": "sku1",
				"variant_type": "size",
				"attribute": "Red"
			}`),
			wantErr: false,
			want: AddVariantOpts{
				SKU:         "sku1",
				Attribute:   "Red",
				VariantType: "size",
			},
		},
		{
			name: "[Error] Empty SKU",
			json: string(`{
				"variant_type": "size",
				"sku": "",
				"attribute": "red"
			}`),
			wantErr: true,
			err:     []string{"sku is a required field"},
		},
		{
			name: "[Error] No SKU",
			json: string(`{
				"variant_type": "size",
				"attribute": "red"
			}`),
			wantErr: true,
			err:     []string{"sku is a required field"},
		},
		{
			name: "[Error] No Attribute",
			json: string(`{
				"variant_type": "size",
				"sku": "red-1"
			}`),
			wantErr: true,
			err:     []string{"attribute is a required field"},
		},
		{
			name: "[Error] Empty Attribute",
			json: string(`{
				"variant_type": "size",
				"sku": "red-1",
				"attribute": ""
			}`),
			wantErr: true,
			err:     []string{"attribute is a required field"},
		},
		{
			name: "[Error] Empty VariantType",
			json: string(`{
				"variant_type": "",
				"sku": "red-1",
				"attribute": "red"
			}`),
			wantErr: true,
			err:     []string{"variant_type is a required field"},
		},
		{
			name: "[Error] No VariantType",
			json: string(`{
				"sku": "red-1",
				"attribute": "red"
			}`),
			wantErr: true,
			err:     []string{"variant_type is a required field"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc AddVariantOpts
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

func TestEditCatalogOpts(t *testing.T) {
	t.Parallel()
	cID, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2612")
	tv := validator.NewValidation()
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    EditCatalogOpts
	}{
		{
			name: "[Ok]",
			json: string(`{
				"id": "5e8821fe1108c87837ef2612",
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"base_price": 1299,
				"retail_price": 1099
			}`),
			wantErr: false,
			want: EditCatalogOpts{
				ID:          cID,
				Name:        "test",
				CategoryID:  []primitive.ObjectID{cID},
				Description: "test description 1",
				Keywords:    []string{"k1", "k2"},
				HSNCode:     "hsnCode1",
				BasePrice:   1299,
				RetailPrice: 1099,
			},
		},
		{
			name: "[Ok] Without Retail Price",
			json: string(`{
				"id": "5e8821fe1108c87837ef2612",
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"base_price": 1299
			}`),
			wantErr: false,
			want: EditCatalogOpts{
				ID:          cID,
				Name:        "test",
				CategoryID:  []primitive.ObjectID{cID},
				Description: "test description 1",
				Keywords:    []string{"k1", "k2"},
				HSNCode:     "hsnCode1",
				BasePrice:   1299,
			},
		},
		{
			name: "[Ok] Without Base Price",
			json: string(`{
				"id": "5e8821fe1108c87837ef2612",
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"retail_price": 1299
			}`),
			wantErr: false,
			want: EditCatalogOpts{
				ID:          cID,
				Name:        "test",
				CategoryID:  []primitive.ObjectID{cID},
				Description: "test description 1",
				Keywords:    []string{"k1", "k2"},
				HSNCode:     "hsnCode1",
				RetailPrice: 1299,
			},
		},
		{
			name: "[Ok] 0 BasePrice",
			json: string(`{
				"id": "5e8821fe1108c87837ef2612",
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"base_price": 0
			}`),
			wantErr: false,
			want: EditCatalogOpts{
				ID:          cID,
				Name:        "test",
				CategoryID:  []primitive.ObjectID{cID},
				Description: "test description 1",
				Keywords:    []string{"k1", "k2"},
				HSNCode:     "hsnCode1",
				BasePrice:   0,
			},
		},
		{
			name: "[Ok] 0 Retail Price",
			json: string(`{
				"id": "5e8821fe1108c87837ef2612",
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"retail_price": 0
			}`),
			wantErr: false,
			want: EditCatalogOpts{
				ID:          cID,
				Name:        "test",
				CategoryID:  []primitive.ObjectID{cID},
				Description: "test description 1",
				Keywords:    []string{"k1", "k2"},
				HSNCode:     "hsnCode1",
				RetailPrice: 0,
			},
		},
		{
			name: "[Ok] Valid BasePrice 0 Retail Price",
			json: string(`{
				"id": "5e8821fe1108c87837ef2612",
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"retail_price": 0,
				"base_price": 1200
			}`),
			wantErr: false,
			want: EditCatalogOpts{
				ID:          cID,
				Name:        "test",
				CategoryID:  []primitive.ObjectID{cID},
				Description: "test description 1",
				Keywords:    []string{"k1", "k2"},
				HSNCode:     "hsnCode1",
				RetailPrice: 0,
				BasePrice:   1200,
			},
		},
		{
			name: "[Error] Base Price Less Than Retail Price",
			json: string(`{
				"id": "5e8821fe1108c87837ef2612",
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"retail_price": 1400,
				"base_price": 1200
			}`),
			wantErr: true,
			err:     []string{"Key: 'EditCatalogOpts.base_price' Error:Field validation for 'base_price' failed on the 'isdefault|gtfield=RetailPrice' tag"},
		},
		{
			name: "[Ok] With Specs",
			json: string(`{
				"id": "5e8821fe1108c87837ef2612",
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"base_price": 1299,
				"retail_price": 1099,
				"specifications": [{
					"Name": "k1",
					"Value": "v1"
				},{
					"Name": "k2",
					"Value": "v2"
				}]
			}`),
			wantErr: false,
			want: EditCatalogOpts{
				ID:          cID,
				Name:        "test",
				CategoryID:  []primitive.ObjectID{cID},
				Description: "test description 1",
				Keywords:    []string{"k1", "k2"},
				HSNCode:     "hsnCode1",
				BasePrice:   1299,
				RetailPrice: 1099,
				Specifications: []specsOpts{
					{
						Name:  "k1",
						Value: "v1",
					},
					{
						Name:  "k2",
						Value: "v2",
					},
				},
			},
		},
		{
			name: "[Ok] With Filter Attribute",
			json: string(`{
				"id": "5e8821fe1108c87837ef2612",
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"base_price": 1299,
				"retail_price": 1099,
				"filter_attr": [{
					"Name": "k1",
					"Value": "v1"
				},{
					"Name": "k2",
					"Value": "v2"
				}]
			}`),
			wantErr: false,
			want: EditCatalogOpts{
				ID:          cID,
				Name:        "test",
				CategoryID:  []primitive.ObjectID{cID},
				Description: "test description 1",
				Keywords:    []string{"k1", "k2"},
				HSNCode:     "hsnCode1",
				BasePrice:   1299,
				RetailPrice: 1099,
				FilterAttribute: []FilterAttribute{
					{
						Name:  "k1",
						Value: "v1",
					},
					{
						Name:  "k2",
						Value: "v2",
					},
				},
			},
		},
		{
			name: "[Error] Passing duplicate keywords",
			json: string(`{
				"id": "5e8821fe1108c87837ef2612",
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"brand_id": "5e8821fe1108c87837ef2611",
				"keywords":  ["k1", "k1"],
				"hsn_code": "hsnCode1",
				"base_price": 1299,
				"retail_price": 1099
			}`),
			wantErr: true,
			err:     []string{"keywords must contain unique values"},
		},
		{
			name: "[Ok] With ETA",
			json: string(`{
				"id": "5e8821fe1108c87837ef2612",
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"brand_id": "5e8821fe1108c87837ef2611",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"base_price": 1299,
				"retail_price": 1099,
				"eta": {
					"min": 1,
					"max": 7,
					"unit": "day"
				}
			}`),
			wantErr: false,
			want: EditCatalogOpts{
				ID:          cID,
				Name:        "test",
				CategoryID:  []primitive.ObjectID{cID},
				Description: "test description 1",
				Keywords:    []string{"k1", "k2"},
				HSNCode:     "hsnCode1",
				BasePrice:   1299,
				RetailPrice: 1099,
				ETA: &etaOpts{
					Min:  1,
					Max:  7,
					Unit: "day",
				},
			},
		},
		{
			name: "[Error] With Invalid ETA Unit",
			json: string(`{
				"id": "5e8821fe1108c87837ef2612",
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"brand_id": "5e8821fe1108c87837ef2611",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"base_price": 1299,
				"retail_price": 1099,
				"eta": {
					"min": 1,
					"max": 7,
					"unit": "year"
				}
			}`),
			wantErr: true,
			err:     []string{"unit must be one of [hour day month]"},
		},
		{
			name: "[Error] With empty Name[1] field specification",
			json: string(`{
				"id": "5e8821fe1108c87837ef2612",
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"brand_id": "5e8821fe1108c87837ef2611",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"specifications": [{
					"Name": "",
					"Value": "v2"
				}],
				"base_price": 1299,
				"retail_price": 1099
			}`),
			wantErr: true,
			err:     []string{"name is a required field"},
		},
		{
			name: "[Error] With empty Value[0] field specification",
			json: string(`{
				"id": "5e8821fe1108c87837ef2612",
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"brand_id": "5e8821fe1108c87837ef2611",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"specifications": [{
					"Name": "k1",
					"Value": ""
				},{
					"Name": "k2",
					"Value": "v2"
				}],
				"base_price": 1299,
				"retail_price": 1099
			}`),
			wantErr: true,
			err:     []string{"value is a required field"},
		},
		{
			name: "[Error] Without ID",
			json: string(`{
				"name": "test",
				"category_id": ["5e8821fe1108c87837ef2612"],
				"description": "test description 1",
				"keywords":  ["k1", "k2"],
				"hsn_code": "hsnCode1",
				"base_price": 1299,
				"retail_price": 1099
			}`),
			wantErr: true,
			err:     []string{"id is a required field"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc EditCatalogOpts
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
func TestKeeperSearchCatalog(t *testing.T) {
	t.Parallel()
	// cID, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2612")
	tv := validator.NewValidation()
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    KeeperSearchCatalogOpts
	}{
		{
			name: "[Ok]",
			json: string(`{
				"name": "test",
				"page": 0
			}`),
			wantErr: false,
			want: KeeperSearchCatalogOpts{
				Name: "test",
				Page: 0,
			},
		},

		{
			name: "[Error] Page Less Than 0",
			json: string(`{
				"name": "test",
				"page": -2
			}`),
			wantErr: true,
			err:     []string{"page must be 0 or greater"},
		},
		{
			name: "[Error] name is required field",
			json: string(`{
				"page": 0
			}`),
			wantErr: true,
			err:     []string{"name is a required field"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc KeeperSearchCatalogOpts
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
func TestKeeperCatalogImpl_DeleteVariant(t *testing.T) {
	t.Parallel()
	cID, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2612")
	vID, _ := primitive.ObjectIDFromHex("603378cb6c45d2a044f167a8")
	tv := validator.NewValidation()
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    DeleteVariantOpts
	}{
		{
			name: "[Ok]",
			json: string(`{
				"catalog_id": "5e8821fe1108c87837ef2612",
				"variant_id": "603378cb6c45d2a044f167a8"
			}`),
			wantErr: false,
			want: DeleteVariantOpts{
				CatalogID: cID,
				VariantID: vID,
			},
		},

		{
			name: "[Error] Catalog ID Missing",
			json: string(`{
				"variant_id": "603378cb6c45d2a044f167a8"
			}`),
			wantErr: true,
			err:     []string{"catalog_id is a required field"},
		},
		{
			name: "[Error] Variant ID Missing",
			json: string(`{
				"catalog_id": "5e8821fe1108c87837ef2612"
			}`),
			wantErr: true,
			err:     []string{"variant_id is a required field"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc DeleteVariantOpts
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

func TestKeeperCatalogImpl_UpdateCatalogStatus(t *testing.T) {
	t.Parallel()
	cID, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2612")
	status := []string{"publish", "unlist", "draft", "archive", "fake"}
	tv := validator.NewValidation()
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    UpdateCatalogStatusOpts
	}{
		{
			name: "[OK] Publish",
			json: string(`{
				"catalog_id":"5e8821fe1108c87837ef2612",
				"status":"publish"
				}`),
			wantErr: false,
			want: UpdateCatalogStatusOpts{
				CatalogID: cID,
				Status:    status[0],
			},
		},
		{
			name: "[OK] unlist",
			json: string(`{
				"catalog_id":"5e8821fe1108c87837ef2612",
				"status":"unlist"
				}`),
			wantErr: false,
			want: UpdateCatalogStatusOpts{
				CatalogID: cID,
				Status:    status[1],
			},
		},
		{
			name: "[OK] draft",
			json: string(`{
				"catalog_id":"5e8821fe1108c87837ef2612",
				"status":"draft"
				}`),
			wantErr: false,
			want: UpdateCatalogStatusOpts{
				CatalogID: cID,
				Status:    status[2],
			},
		},
		{
			name: "[OK] archive",
			json: string(`{
				"catalog_id":"5e8821fe1108c87837ef2612",
				"status":"archive"
				}`),
			wantErr: false,
			want: UpdateCatalogStatusOpts{
				CatalogID: cID,
				Status:    status[3],
			},
		},
		{
			name: "[Error] fake",
			json: string(`{
				"catalog_id":"5e8821fe1108c87837ef2612",
				"status":"fake"
				}`),
			wantErr: true,
			err:     []string{"status must be one of [publish unlist draft archive]"},
		},
		{
			name: "[Error] catalog_id is missing",
			json: string(`{
				"status":"publish"
				}`),
			wantErr: true,
			err:     []string{"catalog_id is a required field"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc UpdateCatalogStatusOpts
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

func TestKeeperCatalogImpl_AddCatalogContent(t *testing.T) {
	t.Parallel()
	cID, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2612")
	bID, _ := primitive.ObjectIDFromHex("603378cb6c45d2a044f167a8")
	fName := "fake file"
	label := ContentLabel{
		Interests: []string{"A", "B"},
		AgeGroup:  []string{"25-30"},
		Gender:    []string{"M", "F"},
	}
	tv := validator.NewValidation()
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    AddCatalogContentOpts
	}{
		{
			name: "[OK]",
			json: string(`{
				"catalog_id":"5e8821fe1108c87837ef2612",
				"brand_id":"603378cb6c45d2a044f167a8",
				"filename":"fake file",
				"label":{
					"interests":["A","B"],
					"age_group":["25-30"],
					"gender":["M","F"]
					}
				}`),
			wantErr: false,
			want: AddCatalogContentOpts{
				BrandID:   bID,
				CatalogID: cID,
				FileName:  fName,
				Label:     &label,
			},
		},
		{
			name: "[Error] catalog_id missing",
			json: string(`{
				"brand_id":"603378cb6c45d2a044f167a8",
				"filename":"fake file",
				"label":{
					"interests":["A","B"],
					"age_group":["25-30"],
					"gender":["M","F"]
					}
				}`),
			wantErr: true,
			err:     []string{"catalog_id is a required field"},
		},
		{
			name: "[Error] Brand is Missing",
			json: string(`{
				"catalog_id":"5e8821fe1108c87837ef2612",
				"filename":"fake file",
				"label":{
					"interests":["A","B"],
					"age_group":["25-30"],
					"gender":["M","F"]
					}
				}`),
			wantErr: true,
			err:     []string{"brand_id is a required field"},
		},
		{
			name: "[Error] File name is missing",
			json: string(`{
				"catalog_id":"5e8821fe1108c87837ef2612",
				"brand_id":"603378cb6c45d2a044f167a8",
				"label":{
					"interests":["A","B"],
					"age_group":["25-30"],
					"gender":["M","F"]
					}
				}`),
			wantErr: true,
			err:     []string{"filename is a required field"},
		},
		{
			name: "[Error] interests is missing",
			json: string(`{
				"catalog_id":"5e8821fe1108c87837ef2612",
				"brand_id":"603378cb6c45d2a044f167a8",
				"filename":"fake file",
				"label":{
					"age_group":["25-30"],
					"gender":["M","F"]
					}
				}`),
			wantErr: true,
			err:     []string{"interests is a required field"},
		},
		{
			name: "[Error] Age Group is missing",
			json: string(`{
				"catalog_id":"5e8821fe1108c87837ef2612",
				"brand_id":"603378cb6c45d2a044f167a8",
				"filename":"fake file",
				"label":{
					"interests":["A","B"],
					"gender":["M","F"]
					}
				}`),
			wantErr: true,
			err:     []string{"age_group is a required field"},
		},
		{
			name: "[Error] Gender is missing",
			json: string(`{
				"catalog_id":"5e8821fe1108c87837ef2612",
				"brand_id":"603378cb6c45d2a044f167a8",
				"filename":"fake file",
				"label":{
					"interests":["A","B"],
					"age_group":["25-30"]
					}
				}`),
			wantErr: true,
			err:     []string{"gender is a required field"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc AddCatalogContentOpts
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

func TestKeeperCatalogImpl_AddCatalogContentImage(t *testing.T) {
	t.Parallel()
	mediaID, _ := primitive.ObjectIDFromHex("603378cb6c45d2a044f167a8")
	label := ContentLabel{
		Interests: []string{"A", "B"},
		AgeGroup:  []string{"25-30"},
		Gender:    []string{"M", "F"},
	}
	cID, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2612")
	tv := validator.NewValidation()
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    AddCatalogContentImageOpts
	}{
		{
			name: "[OK]",
			json: string(`{
				"catalog_id":"5e8821fe1108c87837ef2612",
				"media_id":"603378cb6c45d2a044f167a8",
				"label":{
					"interests":["A","B"],
					"age_group":["25-30"],
					"gender":["M","F"]
					}
				}`),
			wantErr: false,
			want: AddCatalogContentImageOpts{
				CatalogID: cID,
				MediaID:   mediaID,
				Label:     &label,
			},
		},
		{
			name: "[Error] catalog_id missing",
			json: string(`{
				"media_id":"603378cb6c45d2a044f167a8",
				"label":{
					"interests":["A","B"],
					"age_group":["25-30"],
					"gender":["M","F"]
					}
				}`),
			wantErr: true,
			err:     []string{"catalog_id is a required field"},
		},
		{
			name: "[Error] Media ID is Missing",
			json: string(`{
				"catalog_id":"5e8821fe1108c87837ef2612",
				"label":{
					"interests":["A","B"],
					"age_group":["25-30"],
					"gender":["M","F"]
					}
				}`),
			wantErr: true,
			err:     []string{"media_id is a required field"},
		},
		{
			name: "[Error] interests is missing",
			json: string(`{
				"catalog_id":"5e8821fe1108c87837ef2612",
				"media_id":"603378cb6c45d2a044f167a8",
				"label":{
					"age_group":["25-30"],
					"gender":["M","F"]
					}
				}`),
			wantErr: true,
			err:     []string{"interests is a required field"},
		},
		{
			name: "[Error] Age Group is missing",
			json: string(`{
				"catalog_id":"5e8821fe1108c87837ef2612",
				"media_id":"603378cb6c45d2a044f167a8",
				"filename":"fake file",
				"label":{
					"interests":["A","B"],
					"gender":["M","F"]
					}
				}`),
			wantErr: true,
			err:     []string{"age_group is a required field"},
		},
		{
			name: "[Error] Gender is missing",
			json: string(`{
				"catalog_id":"5e8821fe1108c87837ef2612",
				"media_id":"603378cb6c45d2a044f167a8",
				"filename":"fake file",
				"label":{
					"interests":["A","B"],
					"age_group":["25-30"]
					}
				}`),
			wantErr: true,
			err:     []string{"gender is a required field"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc AddCatalogContentImageOpts
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
