package schema

import (
	"encoding/json"
	"go-app/server/validator"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateLiveStreamOpts(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()
	t1, _ := time.Parse(time.RFC3339, "2021-02-14T00:00:00+00:00")
	id1, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2611")
	id2, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2612")
	id3, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2613")
	id4, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2614")
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    CreateLiveStreamOpts
	}{
		{
			name: "[Ok]",
			json: string(`{
				"name": "test stream",
				"featured_image": {
					"src": "https://example.com/test.png"
				},
				"stream_end_image": {
					"src": "https://example.com/test2.png"
				},
				"scheduled_at": "2021-02-14T00:00:00+00:00",
				"influencer_ids": [
					"5e8821fe1108c87837ef2611"
				],
				"catalog_ids": [
					"5e8821fe1108c87837ef2612",
					"5e8821fe1108c87837ef2613",
					"5e8821fe1108c87837ef2614"
				]
			}`),
			wantErr: false,
			want: CreateLiveStreamOpts{
				Name: "test stream",
				FeaturedImage: &Img{
					SRC: "https://example.com/test.png",
				},
				StreamEndImage: &Img{
					SRC: "https://example.com/test2.png",
				},
				ScheduledAt:   t1,
				InfluencerIDs: []primitive.ObjectID{id1},
				CatalogIDs:    []primitive.ObjectID{id2, id3, id4},
			},
		},
		{
			name: "[Error] Without Influencer ID",
			json: string(`{
				"name": "test stream",
				"featured_image": {
					"src": "https://example.com/test.png"
				},
				"stream_end_image": {
					"src": "https://example.com/test2.png"
				},
				"scheduled_at": "2021-02-14T00:00:00+00:00",
				"catalog_ids": [
					"5e8821fe1108c87837ef2612",
					"5e8821fe1108c87837ef2613",
					"5e8821fe1108c87837ef2614"
				]
			}`),
			wantErr: true,
			err:     []string{"influencer_ids is a required field"},
		},
		{
			name: "[Error] Empty Influencer ID Array",
			json: string(`{
				"name": "test stream",
				"featured_image": {
					"src": "https://example.com/test.png"
				},
				"stream_end_image": {
					"src": "https://example.com/test2.png"
				},
				"scheduled_at": "2021-02-14T00:00:00+00:00",
				"catalog_ids": [
					"5e8821fe1108c87837ef2612",
					"5e8821fe1108c87837ef2613",
					"5e8821fe1108c87837ef2614"
				],
				"influencer_ids": []
			}`),
			wantErr: true,
			err:     []string{"influencer_ids must contain at least 1 item"},
		},
		{
			name: "[Error] Without Catalog ID",
			json: string(`{
				"name": "test stream",
				"featured_image": {
					"src": "https://example.com/test.png"
				},
				"stream_end_image": {
					"src": "https://example.com/test2.png"
				},
				"scheduled_at": "2021-02-14T00:00:00+00:00",
				"influencer_ids": [
					"5e8821fe1108c87837ef2612",
					"5e8821fe1108c87837ef2613",
					"5e8821fe1108c87837ef2614"
				]
			}`),
			wantErr: true,
			err:     []string{"catalog_ids is a required field"},
		},
		{
			name: "[Error] Empty Catalog ID Array",
			json: string(`{
				"name": "test stream",
				"featured_image": {
					"src": "https://example.com/test.png"
				},
				"stream_end_image": {
					"src": "https://example.com/test2.png"
				},
				"scheduled_at": "2021-02-14T00:00:00+00:00",
				"influencer_ids": [
					"5e8821fe1108c87837ef2612",
					"5e8821fe1108c87837ef2613",
					"5e8821fe1108c87837ef2614"
				],
				"catalog_ids": []
			}`),
			wantErr: true,
			err:     []string{"catalog_ids must contain at least 1 item"},
		},
		{
			name: "[Error] Without ScheduledAt",
			json: string(`{
				"name": "test stream",
				"featured_image": {
					"src": "https://example.com/test.png"
				},
				"stream_end_image": {
					"src": "https://example.com/test2.png"
				},
				"catalog_ids": [
					"5e8821fe1108c87837ef2612",
					"5e8821fe1108c87837ef2613",
					"5e8821fe1108c87837ef2614"
				],
				"influencer_ids": [
					"5e8821fe1108c87837ef2611"
				]
			}`),
			wantErr: true,
			err:     []string{"scheduled_at is a required field"},
		},
		{
			name: "[Error] Without FeaturedImage",
			json: string(`{
				"name": "test stream",
				"stream_end_image": {
					"src": "https://example.com/test2.png"
				},
				"scheduled_at": "2021-02-14T00:00:00+00:00",
				"catalog_ids": [
					"5e8821fe1108c87837ef2612",
					"5e8821fe1108c87837ef2613",
					"5e8821fe1108c87837ef2614"
				],
				"influencer_ids": [
					"5e8821fe1108c87837ef2611"
				]
			}`),
			wantErr: true,
			err:     []string{"featured_image is a required field"},
		},
		{
			name: "[Error] Without FeaturedImage SRC",
			json: string(`{
				"name": "test stream",
				"featured_image": {},
				"stream_end_image": {
					"src": "https://example.com/test2.png"
				},
				"scheduled_at": "2021-02-14T00:00:00+00:00",
				"catalog_ids": [
					"5e8821fe1108c87837ef2612",
					"5e8821fe1108c87837ef2613",
					"5e8821fe1108c87837ef2614"
				],
				"influencer_ids": [
					"5e8821fe1108c87837ef2611"
				]
			}`),
			wantErr: true,
			err:     []string{"src is a required field"},
		},
		{
			name: "[Error] Without StreamEndImage SRC",
			json: string(`{
				"name": "test stream",
				"featured_image": {
					"src": "https://example.com/test2.png"
				},
				"scheduled_at": "2021-02-14T00:00:00+00:00",
				"catalog_ids": [
					"5e8821fe1108c87837ef2612",
					"5e8821fe1108c87837ef2613",
					"5e8821fe1108c87837ef2614"
				],
				"influencer_ids": [
					"5e8821fe1108c87837ef2611"
				]
			}`),
			wantErr: true,
			err:     []string{"stream_end_image is a required field"},
		},
		{
			name: "[Error] Without name",
			json: string(`{
				"featured_image": {
					"src": "https://example.com/test.png"
				},
				"stream_end_image": {
					"src": "https://example.com/test2.png"
				},
				"scheduled_at": "2021-02-14T00:00:00+00:00",
				"influencer_ids": [
					"5e8821fe1108c87837ef2611"
				],
				"catalog_ids": [
					"5e8821fe1108c87837ef2612",
					"5e8821fe1108c87837ef2613",
					"5e8821fe1108c87837ef2614"
				]
			}`),
			wantErr: true,
			err:     []string{"name is a required field"},
		},
		{
			name: "[Error] Empty name",
			json: string(`{
				"name": "",
				"featured_image": {
					"src": "https://example.com/test.png"
				},
				"stream_end_image": {
					"src": "https://example.com/test2.png"
				},
				"scheduled_at": "2021-02-14T00:00:00+00:00",
				"influencer_ids": [
					"5e8821fe1108c87837ef2611"
				],
				"catalog_ids": [
					"5e8821fe1108c87837ef2612",
					"5e8821fe1108c87837ef2613",
					"5e8821fe1108c87837ef2614"
				]
			}`),
			wantErr: true,
			err:     []string{"name is a required field"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc CreateLiveStreamOpts
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
