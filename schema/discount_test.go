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

func TestCreateDiscountOpts(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()

	cID, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2612")
	vID1, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2613")
	sID, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2614")

	str1 := "2020-02-08T15:04:05+00:00"
	validAfter, _ := time.Parse(time.RFC3339, str1)
	str2 := "2020-02-16T15:04:05+00:00"
	validBefore, _ := time.Parse(time.RFC3339, str2)

	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    CreateDiscountOpts
	}{
		{
			name: "[Ok] Flat Off",
			json: string(`{
				"catalog_id": "5e8821fe1108c87837ef2612",
				"variants_id": ["5e8821fe1108c87837ef2613"],
				"sale_id": "5e8821fe1108c87837ef2614",
				"type": "flat_off",
				"value": 1000
			}`),
			want: CreateDiscountOpts{
				CatalogID:  cID,
				VariantsID: []primitive.ObjectID{vID1},
				SaleID:     sID,
				Type:       model.FlatOffType,
				Value:      1000,
			},
		},
		{
			name: "[Error] Flat Off With MaxValue",
			json: string(`{
				"catalog_id": "5e8821fe1108c87837ef2612",
				"variants_id": ["5e8821fe1108c87837ef2613"],
				"sale_id": "5e8821fe1108c87837ef2614",
				"type": "flat_off",
				"max_value": 10000,
				"value": 1000
			}`),
			wantErr: true,
			err:     []string{"Key: 'CreateDiscountOpts.max_value' Error:Field validation for 'max_value' failed on the 'required_if' tag"},
		},
		{
			name: "[Ok] Percen Off With MaxValue",
			json: string(`{
				"catalog_id": "5e8821fe1108c87837ef2612",
				"variants_id": ["5e8821fe1108c87837ef2613"],
				"sale_id": "5e8821fe1108c87837ef2614",
				"type": "percent_off",
				"max_value": 10000,
				"value": 1000
			}`),
			want: CreateDiscountOpts{
				CatalogID:  cID,
				VariantsID: []primitive.ObjectID{vID1},
				SaleID:     sID,
				Type:       model.PercentOffType,
				Value:      1000,
				MaxValue:   10000,
			},
		},
		{
			name: "[Ok] Percent Off Without MaxValue",
			json: string(`{
				"catalog_id": "5e8821fe1108c87837ef2612",
				"variants_id": ["5e8821fe1108c87837ef2613"],
				"sale_id": "5e8821fe1108c87837ef2614",
				"type": "percent_off",
				"value": 1000
			}`),
			want: CreateDiscountOpts{
				CatalogID:  cID,
				VariantsID: []primitive.ObjectID{vID1},
				SaleID:     sID,
				Type:       model.PercentOffType,
				Value:      1000,
			},
		},
		{
			name: "[Ok] Without SaleID and ValidAfter & ValidBefore",
			json: string(`{
				"catalog_id": "5e8821fe1108c87837ef2612",
				"variants_id": ["5e8821fe1108c87837ef2613"],
				"type": "percent_off",
				"value": 1000,
				"valid_after": "2020-02-08T15:04:05+00:00",
				"valid_before": "2020-02-16T15:04:05+00:00"
			}`),
			want: CreateDiscountOpts{
				CatalogID:   cID,
				VariantsID:  []primitive.ObjectID{vID1},
				Type:        model.PercentOffType,
				Value:       1000,
				ValidAfter:  validAfter,
				ValidBefore: validBefore,
			},
		},
		{
			name: "[Error] With SaleID and ValidAfter",
			json: string(`{
				"catalog_id": "5e8821fe1108c87837ef2612",
				"variants_id": ["5e8821fe1108c87837ef2613"],
				"type": "percent_off",
				"value": 1000,
				"valid_before": "2020-02-16T15:04:05+00:00"
			}`),
			wantErr: true,
			err:     []string{"Key: 'CreateDiscountOpts.valid_after' Error:Field validation for 'valid_after' failed on the 'required_without' tag"},
		},
		{
			name: "[Error] With SaleID and ValidBefore",
			json: string(`{
				"catalog_id": "5e8821fe1108c87837ef2612",
				"variants_id": ["5e8821fe1108c87837ef2613"],
				"type": "percent_off",
				"value": 1000,
				"valid_before": "2020-02-16T15:04:05+00:00"
			}`),
			wantErr: true,
			err:     []string{"Key: 'CreateDiscountOpts.valid_after' Error:Field validation for 'valid_after' failed on the 'required_without' tag"},
		},
		{
			name: "[Error] ValidAfter greater than ValidBefore",
			json: string(`{
				"catalog_id": "5e8821fe1108c87837ef2612",
				"variants_id": ["5e8821fe1108c87837ef2613"],
				"type": "percent_off",
				"value": 1000,
				"valid_before": "2020-02-16T15:04:05+00:00",
				"valid_after": "2020-02-18T15:04:05+00:00"
			}`),
			wantErr: true,
			err:     []string{"Key: 'CreateDiscountOpts.valid_before' Error:Field validation for 'valid_before' failed on the 'isdefault|gtfield=ValidAfter' tag"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc CreateDiscountOpts
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

func TestCreateSaleOpts(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()

	str1 := "2020-02-08T15:04:05+00:00"
	validAfter, _ := time.Parse(time.RFC3339, str1)
	str2 := "2020-02-16T15:04:05+00:00"
	validBefore, _ := time.Parse(time.RFC3339, str2)

	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    CreateSaleOpts
	}{
		{
			name: "[Ok]",
			json: string(`{
				"name": "test sale",
				"banner": {
					"src": "https://img.com/x.png"
				},
				"valid_before": "2020-02-16T15:04:05+00:00",
				"valid_after": "2020-02-08T15:04:05+00:00"
			}`),
			wantErr: false,
			want: CreateSaleOpts{
				Name: "test sale",
				Banner: Img{
					SRC: "https://img.com/x.png",
				},
				ValidBefore: validBefore,
				ValidAfter:  validAfter,
			},
		},
		{
			name: "[Error] Without Banner",
			json: string(`{
				"name": "test sale",
				"valid_before": "2020-02-16T15:04:05+00:00",
				"valid_after": "2020-02-08T15:04:05+00:00"
			}`),
			wantErr: true,
			err:     []string{"src is a required field"},
		},
		{
			name: "[Error] Without Name",
			json: string(`{
				"banner": {
					"src": "https://img.com/x.png"
				},
				"valid_before": "2020-02-16T15:04:05+00:00",
				"valid_after": "2020-02-08T15:04:05+00:00"
			}`),
			wantErr: true,
			err:     []string{"name is a required field"},
		},
		{
			name: "[Error] Without ValidAfter",
			json: string(`{
				"banner": {
					"src": "https://img.com/x.png"
				},
				"valid_before": "2020-02-16T15:04:05+00:00",
				"name": "test sale"
			}`),
			wantErr: true,
			err:     []string{"valid_after is a required field"},
		},
		{
			name: "[Error] Without ValidBefore",
			json: string(`{
				"banner": {
					"src": "https://img.com/x.png"
				},
				"valid_after": "2020-02-16T15:04:05+00:00",
				"name": "test sale"
			}`),
			wantErr: true,
			err:     []string{"valid_before is a required field"},
		},
		{
			name: "[Error] ValidAfter greater than ValidBefore",
			json: string(`{
				"banner": {
					"src": "https://img.com/x.png"
				},
				"valid_after": "2020-02-18T15:04:05+00:00",
				"valid_before": "2020-02-16T15:04:05+00:00",
				"name": "test sale"
			}`),
			wantErr: true,
			err:     []string{"valid_before must be greater than ValidAfter"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc CreateSaleOpts
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
