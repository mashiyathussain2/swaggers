package app

import (
	"context"
	"fmt"
	"go-app/mock"
	"go-app/model"
	"go-app/schema"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestInventoryImpl_CreateInventory(t *testing.T) {
	t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	catalog := schema.GetCatalogResp{
		ID: primitive.NewObjectID(),
		Variants: []schema.VariantInfo{
			{
				ID: primitive.NewObjectID(),
			},
		},
	}

	inventory := model.Inventory{
		ID:          primitive.NewObjectID(),
		CatalogID:   primitive.NewObjectID(),
		VariantID:   primitive.NewObjectID(),
		UnitInStock: 10,
	}
	catWithInv := schema.GetCatalogResp{
		ID: inventory.CatalogID,
		Variants: []schema.VariantInfo{
			{
				ID: inventory.VariantID,
			},
		},
	}
	inventoryDB := app.MongoDB.Client.Database(app.Config.InventoryConfig.DBName)
	inventoryDB.Collection(model.InventoryColl).InsertOne(context.TODO(), inventory)
	type args struct {
		opts *schema.CreateInventoryOpts
	}
	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(tt *TC, kc *mock.MockKeeperCatalog)
		// validator  func(*testing.T, *schema.CreateCatalogOpts, *schema.CreateCatalogResp)
	}
	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {
			},
			prepare: func(tt *TC) {
				opts := schema.CreateInventoryOpts{
					CatalogID: catalog.ID,
					VariantID: catalog.Variants[0].ID,
					Unit:      5,
				}
				tt.args.opts = &opts
			},
			wantErr: false,
		},
		{
			name: "[Error] Inventory already exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				opts := schema.CreateInventoryOpts{
					CatalogID: catWithInv.ID,
					VariantID: catWithInv.Variants[0].ID,
					Unit:      5,
				}
				tt.args.opts = &opts
				tt.err = errors.Errorf("inventory already exist")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ii := &InventoryImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockKeeperCatalog)
			tt.prepare(&tt)

			res, err := ii.CreateInventory(tt.args.opts)
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, res)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.True(t, res.IsZero())
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestInventoryImpl_UpdateInventory(t *testing.T) {
	// t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	inventory := model.Inventory{
		ID:          primitive.NewObjectID(),
		CatalogID:   primitive.NewObjectID(),
		VariantID:   primitive.NewObjectID(),
		UnitInStock: 0,
		Status: &model.InventoryStatus{
			Value: model.OutOfStockStatus,
		},
	}
	inventoryDB := app.MongoDB.Client.Database(app.Config.InventoryConfig.DBName)
	inventoryDB.Collection(model.InventoryColl).InsertOne(context.TODO(), inventory)

	type args struct {
		opts *schema.UpdateInventoryOpts
	}
	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(tt *TC, kc *mock.MockKeeperCatalog)
		validator  func(*testing.T, *schema.UpdateInventoryOpts)
	}
	tests := []TC{
		{
			name: "[Ok] set 0 to +ve",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				opts := schema.UpdateInventoryOpts{
					ID: inventory.ID,
					Operation: &schema.UpdateInventoryOperation{
						Operator: "set",
						Unit:     2,
					},
				}
				tt.args.opts = &opts
			},
			wantErr: false,
			validator: func(t *testing.T, opts *schema.UpdateInventoryOpts) {
				var inv model.Inventory
				inventoryDB.Collection(model.InventoryColl).FindOne(context.TODO(), bson.M{"_id": opts.ID}).Decode(&inv)
				assert.Equal(t, inv.Status.Value, model.InStockStatus)
				assert.WithinDuration(t, time.Now().UTC(), inv.Status.CreatedAt, 100*time.Millisecond)
			},
		},
		{
			name: "[Ok] set +ve to 0",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				opts := schema.UpdateInventoryOpts{
					ID: inventory.ID,
					Operation: &schema.UpdateInventoryOperation{
						Operator: "set",
						Unit:     0,
					},
				}
				tt.args.opts = &opts
			},
			wantErr: false,
			validator: func(t *testing.T, opts *schema.UpdateInventoryOpts) {
				var inv model.Inventory
				inventoryDB.Collection(model.InventoryColl).FindOne(context.TODO(), bson.M{"_id": opts.ID}).Decode(&inv)
				assert.Equal(t, inv.Status.Value, model.OutOfStockStatus)
				assert.WithinDuration(t, time.Now().UTC(), inv.Status.CreatedAt, 100*time.Millisecond)
			},
		},
		{
			name: "[Ok] add to Out of Stock",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				opts := schema.UpdateInventoryOpts{
					ID: inventory.ID,
					Operation: &schema.UpdateInventoryOperation{
						Operator: "add",
						Unit:     5,
					},
				}
				tt.args.opts = &opts
			},
			wantErr: false,
			validator: func(t *testing.T, opts *schema.UpdateInventoryOpts) {
				var inv model.Inventory
				inventoryDB.Collection(model.InventoryColl).FindOne(context.TODO(), bson.M{"_id": opts.ID}).Decode(&inv)
				assert.Equal(t, inv.Status.Value, model.InStockStatus)
				assert.WithinDuration(t, time.Now().UTC(), inv.Status.CreatedAt, 100*time.Millisecond)
			},
		},
		{
			name: "[Ok] add to In Stock",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				inventoryInStock := model.Inventory{
					ID:          primitive.NewObjectID(),
					CatalogID:   primitive.NewObjectID(),
					VariantID:   primitive.NewObjectID(),
					UnitInStock: 5,
					Status: &model.InventoryStatus{
						Value: model.InStockStatus,
					},
				}
				inventoryDB.Collection(model.InventoryColl).InsertOne(context.TODO(), inventoryInStock)

				opts := schema.UpdateInventoryOpts{
					ID: inventoryInStock.ID,
					Operation: &schema.UpdateInventoryOperation{
						Operator: "add",
						Unit:     5,
					},
				}
				tt.args.opts = &opts
			},
			wantErr: false,
			validator: func(t *testing.T, opts *schema.UpdateInventoryOpts) {
				var inv model.Inventory
				inventoryDB.Collection(model.InventoryColl).FindOne(context.TODO(), bson.M{"_id": opts.ID}).Decode(&inv)
				assert.Equal(t, inv.Status.Value, model.InStockStatus)
				// assert.WithinDuration(t, time.Now().UTC(), inv.Status.CreatedAt, 100*time.Millisecond)
				assert.Equal(t, inv.UnitInStock, 10)
			},
		},
		{
			name: "[Ok] Sub complete unit",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				inventoryInStock := model.Inventory{
					ID:          primitive.NewObjectID(),
					CatalogID:   primitive.NewObjectID(),
					VariantID:   primitive.NewObjectID(),
					UnitInStock: 5,
					Status: &model.InventoryStatus{
						Value: model.InStockStatus,
					},
				}
				inventoryDB.Collection(model.InventoryColl).InsertOne(context.TODO(), inventoryInStock)

				opts := schema.UpdateInventoryOpts{
					ID: inventoryInStock.ID,
					Operation: &schema.UpdateInventoryOperation{
						Operator: "subtract",
						Unit:     5,
					},
				}
				tt.args.opts = &opts
			},
			wantErr: false,
			validator: func(t *testing.T, opts *schema.UpdateInventoryOpts) {
				var inv model.Inventory
				inventoryDB.Collection(model.InventoryColl).FindOne(context.TODO(), bson.M{"_id": opts.ID}).Decode(&inv)
				fmt.Println(inv)
				assert.Equal(t, inv.Status.Value, model.OutOfStockStatus)
				assert.WithinDuration(t, time.Now().UTC(), inv.Status.CreatedAt, 100*time.Millisecond)
				assert.Equal(t, inv.UnitInStock, 0)
			},
		},
		{
			name: "[Ok] Sub some Qty",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				inventoryInStock := model.Inventory{
					ID:          primitive.NewObjectID(),
					CatalogID:   primitive.NewObjectID(),
					VariantID:   primitive.NewObjectID(),
					UnitInStock: 5,
					Status: &model.InventoryStatus{
						Value: model.InStockStatus,
					},
				}
				inventoryDB.Collection(model.InventoryColl).InsertOne(context.TODO(), inventoryInStock)

				opts := schema.UpdateInventoryOpts{
					ID: inventoryInStock.ID,
					Operation: &schema.UpdateInventoryOperation{
						Operator: "subtract",
						Unit:     2,
					},
				}
				tt.args.opts = &opts
			},
			wantErr: false,
			validator: func(t *testing.T, opts *schema.UpdateInventoryOpts) {
				var inv model.Inventory
				inventoryDB.Collection(model.InventoryColl).FindOne(context.TODO(), bson.M{"_id": opts.ID}).Decode(&inv)
				assert.Equal(t, inv.Status.Value, model.InStockStatus)
				// assert.WithinDuration(t, time.Now().UTC(), inv.Status.CreatedAt, 100*time.Millisecond)
				assert.Equal(t, inv.UnitInStock, 3)
			},
		},
		{
			name: "[Error] Subtract < 0",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				opts := schema.UpdateInventoryOpts{
					ID: inventory.ID,
					Operation: &schema.UpdateInventoryOperation{
						Operator: "subtract",
						Unit:     20,
					},
				}
				tt.args.opts = &opts
				tt.err = errors.Errorf("inventory for id: %s, cannot be negative", tt.args.opts.ID)

			},
			wantErr: true,
		},
		{
			name: "[Error] Inventory ID doesn't exist",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				opts := schema.UpdateInventoryOpts{
					ID: primitive.NewObjectID(),
					Operation: &schema.UpdateInventoryOperation{
						Operator: "set",
						Unit:     2,
					},
				}
				tt.args.opts = &opts
				tt.err = errors.Errorf("unable to find the inventory with id: %s", opts.ID.Hex())
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ii := &InventoryImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockKeeperCatalog)
			tt.prepare(&tt)

			err := ii.UpdateInventory(tt.args.opts)
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validator(t, tt.args.opts)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestInventoryImpl_SetOutOfStock(t *testing.T) {
	// t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	inventory := model.Inventory{
		ID:          primitive.NewObjectID(),
		CatalogID:   primitive.NewObjectID(),
		VariantID:   primitive.NewObjectID(),
		UnitInStock: 10,
		Status: &model.InventoryStatus{
			Value: model.InStockStatus,
		},
	}
	inventoryDB := app.MongoDB.Client.Database(app.Config.InventoryConfig.DBName)
	inventoryDB.Collection(model.InventoryColl).InsertOne(context.TODO(), inventory)

	type args struct {
		id primitive.ObjectID
	}
	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		err        error
		prepare    func(*TC)
		buildStubs func(tt *TC, kc *mock.MockKeeperCatalog)
		validator  func(*testing.T, primitive.ObjectID)
	}
	tests := []TC{
		{
			name: "[Ok] set in stock to out of stock",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.args.id = inventory.ID
			},
			wantErr: false,
			validator: func(t *testing.T, id primitive.ObjectID) {
				var inv model.Inventory
				inventoryDB.Collection(model.InventoryColl).FindOne(context.TODO(), bson.M{"_id": id}).Decode(&inv)
				assert.Equal(t, inv.Status.Value, model.OutOfStockStatus)
				assert.Zero(t, inv.UnitInStock)
				assert.WithinDuration(t, time.Now().UTC(), inv.Status.CreatedAt, 100*time.Millisecond)
			},
		},

		{
			name: "[Error] Inventory ID doesn't exist ",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.args.id = primitive.NewObjectID()
				tt.err = errors.Errorf("unable to find inventory with id: %s", tt.args.id.Hex())
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ii := &InventoryImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockKeeperCatalog)
			tt.prepare(&tt)

			err := ii.SetOutOfStock(tt.args.id)
			if !tt.wantErr {
				assert.Nil(t, err)
				tt.validator(t, tt.args.id)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}

func TestInventoryImpl_CheckInventoryExists(t *testing.T) {
	// t.Parallel()
	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	inventory := model.Inventory{
		ID:          primitive.NewObjectID(),
		CatalogID:   primitive.NewObjectID(),
		VariantID:   primitive.NewObjectID(),
		UnitInStock: 10,
		Status: &model.InventoryStatus{
			Value: model.InStockStatus,
		},
	}
	inventoryDB := app.MongoDB.Client.Database(app.Config.InventoryConfig.DBName)
	inventoryDB.Collection(model.InventoryColl).InsertOne(context.TODO(), inventory)

	inventoryOS := model.Inventory{
		ID:          primitive.NewObjectID(),
		CatalogID:   primitive.NewObjectID(),
		VariantID:   primitive.NewObjectID(),
		UnitInStock: 0,
		Status: &model.InventoryStatus{
			Value: model.OutOfStockStatus,
		},
	}
	inventoryDB.Collection(model.InventoryColl).InsertOne(context.TODO(), inventoryOS)

	type args struct {
		cat_id primitive.ObjectID
		var_id primitive.ObjectID
		qty    int
	}
	type fields struct {
		App    *App
		DB     *mongo.Database
		Logger *zerolog.Logger
	}
	type TC struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		want       bool
		err        error
		prepare    func(*TC)
		buildStubs func(tt *TC, kc *mock.MockKeeperCatalog)
	}
	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.args.cat_id = inventory.CatalogID
				tt.args.var_id = inventory.VariantID
				tt.args.qty = 5
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "[Ok] with no inventory",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.args.cat_id = inventoryOS.CatalogID
				tt.args.var_id = inventoryOS.VariantID
				tt.args.qty = 5
			},
			wantErr: false,
			want:    false,
		},
		{
			name: "[Ok] with insufficient inventory",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.args.cat_id = inventory.CatalogID
				tt.args.var_id = inventory.VariantID
				tt.args.qty = 15
			},
			wantErr: false,
			want:    false,
		},
		{
			name: "[Error] with -ve inventory",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.args.cat_id = inventory.CatalogID
				tt.args.var_id = inventory.VariantID
				tt.args.qty = -5
				tt.err = errors.Errorf("quantity must be greater than 0")
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "[Error] category ID doesn't exist ",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.args.cat_id = primitive.NewObjectID()
				tt.args.var_id = inventory.VariantID
				tt.args.qty = 5
				tt.err = errors.Errorf("inventory not found")
			},
			wantErr: true,
		},
		{
			name: "[Error] variant ID doesn't exist ",
			fields: fields{
				App:    app,
				DB:     app.MongoDB.Client.Database(app.Config.GroupConfig.DBName),
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, kc *mock.MockKeeperCatalog) {

			},
			prepare: func(tt *TC) {
				tt.args.cat_id = inventory.CatalogID
				tt.args.var_id = primitive.NewObjectID()
				tt.args.qty = 5
				tt.err = errors.Errorf("inventory not found")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ii := &InventoryImpl{
				App:    tt.fields.App,
				DB:     tt.fields.DB,
				Logger: tt.fields.Logger,
			}
			mockKeeperCatalog := mock.NewMockKeeperCatalog(ctrl)
			tt.fields.App.KeeperCatalog = mockKeeperCatalog
			tt.buildStubs(&tt, mockKeeperCatalog)
			tt.prepare(&tt)

			found, err := ii.CheckInventoryExists(tt.args.cat_id, tt.args.var_id, tt.args.qty)
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.Equal(t, found, tt.want)
			}

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
				assert.False(t, found)
			}
		})
	}
}
