package app

import (
	"go-app/schema"
	"testing"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestKeeperCatalogImpl_CreateCatalog(t *testing.T) {
	app := NewTestApp(getTestConfig())
	// defer CleanTestApp(app)
	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type args struct {
		opts *schema.CreateCatalogOpts
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Ok",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.KeeperCatalogConfig.DBName),
				Logger: app.Logger,
			},
			args: args{
				opts: &schema.CreateCatalogOpts{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kc := &KeeperCatalogImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			kc.CreateCatalog(tt.args.opts)
		})
	}
}
