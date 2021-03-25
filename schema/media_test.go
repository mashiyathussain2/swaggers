package schema

import (
	"encoding/json"
	"go-app/server/validator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateVideoUploadTokenOpts(t *testing.T) {
	t.Parallel()
	tv := validator.NewValidation()
	tests := []struct {
		name    string
		json    string
		wantErr bool
		err     []string
		want    GenerateVideoUploadTokenOpts
	}{
		{
			name: "[Ok]",
			json: string(`{
				"file_name": "xyx.png"
			}`),
			wantErr: false,
			want: GenerateVideoUploadTokenOpts{
				FileName: "xyx.png",
			},
		},
		{
			name: "[Error] Empty filename",
			json: string(`{
				"file_name": ""
			}`),
			wantErr: true,
			err:     []string{"file_name is a required field"},
		},
		{
			name: "[Error] No filename",
			json: string(`{
			
			}`),
			wantErr: true,
			err:     []string{"file_name is a required field"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sc GenerateVideoUploadTokenOpts
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
