package app

import (
	"go-app/schema"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCategoryImpl_CreateCategory(t *testing.T) {
	t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	// opts := schema.GetRandomCreateCategoryOpts()

	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.CreateCategoryOpts
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *schema.CreateCategoryResp
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CategoryImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			got, err := c.CreateCategory(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("CategoryImpl.CreateCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CategoryImpl.CreateCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}
