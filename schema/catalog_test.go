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
				"retail_price": 1099
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
				FilterAttribute: []filterAttribute{
					{
						Name:  "Color",
						Value: "Red",
					},
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
				"variants": [
					{
						"sku": "sku1",
						"base_price": 1299,
						"retail_price": 1099
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
						SKU: "sku1",
					},
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
				"variants": [
					{
						"sku": "sku1",
						"base_price": 1299,
						"retail_price": 1099
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
					"retail_price": 1099
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
					"retail_price": 1099
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
					"retail_price": 1099
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
					"retail_price": 1099
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
					"retail_price": 1099
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
					"base_price": 1099
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
					"retail_price": 0
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
					"retail_price": 1099
				}`),
			wantErr: true,
			err:     []string{"category_id must contain more than 0 items"},
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
				"sku": "sku1"
			}`),
			wantErr: false,
			want: CreateVariantOpts{
				SKU: "sku1",
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
				"sku": ""
			}`),
			wantErr: true,
			err:     []string{"sku is a required field"},
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
				assert.Equal(t, tt.want, sc)
			}
		})
	}
}
