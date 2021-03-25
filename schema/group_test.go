package schema

import (
	"encoding/json"
	"go-app/model"
	"go-app/server/validator"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGroupImpl_CreateCatalogGroup(t *testing.T) {
	t.Parallel()
	name := "Test Group"
	cID1, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2612")
	cID2, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2632")
	cIDs := []primitive.ObjectID{cID1, cID2}
	tv := validator.NewValidation()

	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    CreateCatalogGroupOpts
	}{
		{
			name: "[Ok]",
			json: string(`{
				"basis": "Test Group",
				"ids":["5e8821fe1108c87837ef2612", "5e8821fe1108c87837ef2632"]
			}`),
			wantErr: false,
			want: CreateCatalogGroupOpts{
				Basis: name,
				IDs:   cIDs,
			},
		},
		{
			name: "[Error] Basis is Missing",
			json: string(`{
				"ids":["5e8821fe1108c87837ef2612", "5e8821fe1108c87837ef2632"]
			}`),
			wantErr: true,
			err:     []string{"basis is a required field"},
		},
		{
			name: "[Error] IDs is Missing",
			json: string(`{
				"basis":"Test Name"
			}`),
			wantErr: true,
			err:     []string{"ids is a required field"},
		},
		{
			name: "[Error] IDs is Missing",
			json: string(`{
				"basis":"Test Name",
				"ids":[]
			}`),
			wantErr: true,
			err:     []string{"ids must contain at least 1 item"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc CreateCatalogGroupOpts
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

func TestGroupImpl_GetGroups(t *testing.T) {
	page0 := 0
	pagePos := 5
	status := model.Publish
	tv := validator.NewValidation()

	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    GetGroupsOpts
	}{
		{
			name: "[Ok] page +ve",
			json: string(`{
				"page": 5,
				"status":"publish"
			}`),
			wantErr: false,
			want: GetGroupsOpts{
				Page:   pagePos,
				Status: status,
			},
		},
		{
			name: "[Ok] page 0",
			json: string(`{
				"page": 0,
				"status":"publish"
			}`),
			wantErr: false,
			want: GetGroupsOpts{
				Page:   page0,
				Status: status,
			},
		},
		{
			name: "[Error] Page -ve",
			json: string(`{
				"page": -5,
				"status":"publish"
			}`),
			wantErr: true,
			err:     []string{"page must be 0 or greater"},
		},
		{
			name: "[Ok] Page not entered",
			json: string(`{
				"status":"publish"
			}`),
			wantErr: false,
			want: GetGroupsOpts{
				Page:   page0,
				Status: model.Publish,
			},
		},
		{
			name: "[Error] Status not entered",
			json: string(`{
				"page": 5
			}`),
			wantErr: true,
			err:     []string{"status is a required field"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc GetGroupsOpts
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

func TestGroupImpl_GetGroupsByCatalogIDOpts(t *testing.T) {
	tv := validator.NewValidation()
	cID, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2612")
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    GetGroupsByCatalogIDOpts
	}{
		{
			name: "[OK]",
			json: string(`{
				"page": 5,
				"id":"5e8821fe1108c87837ef2612"
			}`),
			wantErr: false,
			want: GetGroupsByCatalogIDOpts{
				Page: 5,
				ID:   cID,
			},
		},
		{
			name: "[Ok] page 0",
			json: string(`{
				"page": 0,
				"id":"5e8821fe1108c87837ef2612"
			}`),
			wantErr: false,
			want: GetGroupsByCatalogIDOpts{
				Page: 0,
				ID:   cID,
			},
		},
		{
			name: "[Error] Page -ve",
			json: string(`{
				"page": -5,
				"id":"5e8821fe1108c87837ef2612"
			}`),
			wantErr: true,
			err:     []string{"page must be 0 or greater"},
		},
		{
			name: "[Ok] Page not entered",
			json: string(`{
				"id":"5e8821fe1108c87837ef2612"
			}`),
			wantErr: false,
			want: GetGroupsByCatalogIDOpts{
				Page: 0,
				ID:   cID,
			},
		},
		{
			name: "[Error] Catalog ID not entered",
			json: string(`{
				"page": 5
			}`),
			wantErr: true,
			err:     []string{"id is a required field"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc GetGroupsByCatalogIDOpts
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

func TestGroupImpl_KeeperGetGroupsByCatalogID(t *testing.T) {
	tv := validator.NewValidation()
	cID, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2612")
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    KeeperGetGroupsByCatalogIDOpts
	}{
		{

			name: "[OK]",
			json: string(`{
					"page": 5,
					"id":"5e8821fe1108c87837ef2612",
					"status":"publish"
				}`),
			wantErr: false,
			want: KeeperGetGroupsByCatalogIDOpts{
				Page:   5,
				ID:     cID,
				Status: model.Publish,
			},
		},
		{
			name: "[OK] without status",
			json: string(`{
					"page": 5,
					"id":"5e8821fe1108c87837ef2612"
				}`),
			wantErr: false,
			want: KeeperGetGroupsByCatalogIDOpts{
				Page: 5,
				ID:   cID,
			},
		},
		{
			name: "[OK] Page 0",
			json: string(`{
					"page": 0,
					"id":"5e8821fe1108c87837ef2612"
				}`),
			wantErr: false,
			want: KeeperGetGroupsByCatalogIDOpts{
				Page: 0,
				ID:   cID,
			},
		},
		{
			name: "[OK] Page Missing",
			json: string(`{
					"id":"5e8821fe1108c87837ef2612"
				}`),
			wantErr: false,
			want: KeeperGetGroupsByCatalogIDOpts{
				Page: 0,
				ID:   cID,
			},
		},
		{
			name: "[Error] Page -ve",
			json: string(`{
					"page":-2,
					"id":"5e8821fe1108c87837ef2612"
				}`),
			wantErr: true,
			err:     []string{"page must be 0 or greater"},
		},
		{
			name: "[Error] Catalog ID missing",
			json: string(`{
					"page":0
				}`),
			wantErr: true,
			err:     []string{"id is a required field"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc KeeperGetGroupsByCatalogIDOpts
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

func TestGroupImpl_AddCatalogsInTheGroup(t *testing.T) {
	tv := validator.NewValidation()
	gID, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2612")
	cID1, _ := primitive.ObjectIDFromHex("6048e7243b421c459b1acf35")
	cID2, _ := primitive.ObjectIDFromHex("6048e759d39654ad4458751f")
	cIDs := []primitive.ObjectID{
		cID1,
		cID2,
	}

	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    AddCatalogsInTheGroupOpts
	}{
		{
			name: "[OK]",
			json: string(`{
					"id":"5e8821fe1108c87837ef2612",
					"catalog_ids":["6048e7243b421c459b1acf35","6048e759d39654ad4458751f"]
				}`),
			wantErr: false,
			want: AddCatalogsInTheGroupOpts{
				ID:         gID,
				CatalogIDs: cIDs,
			},
		},
		{
			name: "[Error] id (group id) missing",
			json: string(`{
					"catalog_ids":["6048e7243b421c459b1acf35","6048e759d39654ad4458751f"]
				}`),
			wantErr: true,
			err:     []string{"id is a required field"},
		},
		{
			name: "[Error] catalog ids  missing",
			json: string(`{
					"id":"6048e7243b421c459b1acf35"
				}`),
			wantErr: true,
			err:     []string{"catalog_ids must contain more than 0 items"},
		},
		{
			name: "[Error] catalog ids empty array",
			json: string(`{
					"id":"6048e7243b421c459b1acf35",
					"catalog_ids":[]
				}`),
			wantErr: true,
			err:     []string{"catalog_ids must contain more than 0 items"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc AddCatalogsInTheGroupOpts
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

func TestGroupImpl_UpdateGroupStatus(t *testing.T) {
	tv := validator.NewValidation()
	gID, _ := primitive.ObjectIDFromHex("5e8821fe1108c87837ef2612")

	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    UpdateGroupStatusOpts
	}{
		{
			name: "[OK]",
			json: string(`{
					"id":"5e8821fe1108c87837ef2612",
					"status":"publish"
				}`),
			wantErr: false,
			want: UpdateGroupStatusOpts{
				ID:     gID,
				Status: model.Publish,
			},
		},
		{
			name: "[OK] Unlist",
			json: string(`{
					"id":"5e8821fe1108c87837ef2612",
					"status":"unlist"
				}`),
			wantErr: false,
			want: UpdateGroupStatusOpts{
				ID:     gID,
				Status: model.Unlist,
			},
		},
		{
			name: "[OK] archive",
			json: string(`{
					"id":"5e8821fe1108c87837ef2612",
					"status":"archive"
				}`),
			wantErr: false,
			want: UpdateGroupStatusOpts{
				ID:     gID,
				Status: model.Archive,
			},
		},
		{
			name: "[Error] id missing",
			json: string(`{
					"status":"publish"
				}`),
			wantErr: true,
			err:     []string{"id is a required field"},
		},
		{
			name: "[Error] status not in oneof",
			json: string(`{
					"id":"5e8821fe1108c87837ef2612",
					"status":"fake"
				}`),
			wantErr: true,
			err:     []string{"status must be one of [publish unlist archive]"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc UpdateGroupStatusOpts
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
