package app

import (
	"context"
	"fmt"
	"go-app/model"
	"go-app/schema"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Group service allows `Group` to execute admin operations.
type Group interface {
	CreateCatalogGroup(*schema.CreateCatalogGroupOpts) (schema.CreateGroupResp, []error)
	GetCatalogsByGroupID(primitive.ObjectID, int) ([]schema.GetCatalogByGroupIDResp, error)
	GetGroups(*schema.GetGroupsOpts) ([]schema.GroupResp, error)
	GetGroupsByCatalogID(*schema.GetGroupsByCatalogIDOpts) ([]schema.GetGroupsByCatalogIDResp, error)
	KeeperGetGroupsByCatalogID(*schema.KeeperGetGroupsByCatalogIDOpts) ([]schema.GroupResp, error)
	AddCatalogsInTheGroup(opts *schema.AddCatalogsInTheGroupOpts) (bool, []error)
	EditGroup(*schema.EditGroupOpts) (*schema.EditGroupResp, []error)
	GetGroupsByCatalogName(string, int) ([]schema.GetGroupsByCatalogNameResp, error)
}

// GroupImpl implements group related operations
type GroupImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// GroupOpts contains arguments required to create a new instance of Group
type GroupOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

//InitGroup returns Group instance
func InitGroup(opts *GroupOpts) Group {
	return &GroupImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
}

// CreateCatalogGroup inserts a new Group with specified catalog ID to DB
func (gi *GroupImpl) CreateCatalogGroup(opts *schema.CreateCatalogGroupOpts) (schema.CreateGroupResp, []error) {

	ctx := context.TODO()

	var errorResp []error
	var createGroupResp schema.CreateGroupResp
	catalogs, err := gi.App.KeeperCatalog.GetCatalogByIDs(ctx, opts.IDs)
	if err != nil {
		return createGroupResp, []error{errors.Wrap(err, "unable to fetch catalogs by ids")}
	}
	catalogMap := make(map[primitive.ObjectID]schema.GetCatalogResp)
	for i := 0; i < len(catalogs); i++ {
		catalogMap[catalogs[i].ID] = catalogs[i]
	}

	if len(catalogs) != len(opts.IDs) {
		for i := 0; i < len(opts.IDs); i++ {
			_, ok := catalogMap[opts.IDs[i]]
			if !ok {
				errorResp = append(errorResp, errors.Errorf("catalog with id: %s not found", opts.IDs[i].Hex()))
			}
		}
		return createGroupResp, errorResp
	}
	status := &model.GroupStatus{
		Value:     model.Unlist,
		CreatedAt: time.Now().UTC(),
	}
	newGroup := model.Group{
		Basis:      opts.Basis,
		CatalogIDs: opts.IDs,
		CreatedAt:  time.Now().UTC(),
		Status:     status,
	}
	res, err := gi.DB.Collection(model.GroupColl).InsertOne(ctx, newGroup)
	if err != nil {
		return createGroupResp, []error{errors.Wrap(err, "unable to add group to database")}
	}
	createGroupResp = schema.CreateGroupResp{
		ID:         res.InsertedID.(primitive.ObjectID),
		Basis:      opts.Basis,
		CatalogIDs: opts.IDs,
		Status:     model.Unlist,
	}
	return createGroupResp, nil
}

// GetCatalogsByGroupID returns a list of Catalogs in the Group with Given group ID
func (gi *GroupImpl) GetCatalogsByGroupID(id primitive.ObjectID, page int) ([]schema.GetCatalogByGroupIDResp, error) {

	matchStage := bson.D{
		{
			Key: "$match", Value: bson.M{
				"_id": id,
			},
		},
	}
	unwindStage := bson.D{
		{
			Key: "$unwind", Value: bson.M{
				"path": "$catalog_ids",
			},
		},
	}
	limitStage := bson.D{
		{Key: "$limit", Value: gi.App.Config.PageSize},
	}
	skipStage := bson.D{
		{Key: "$skip", Value: gi.App.Config.PageSize * page},
	}

	lookupStage := bson.D{
		{
			Key: "$lookup", Value: bson.M{
				"from":         model.CatalogColl,
				"localField":   "catalog_ids",
				"foreignField": "_id",
				"as":           "catalog_info",
			},
		},
	}
	replaceRootStage := bson.D{
		{
			Key: "$replaceRoot", Value: bson.M{
				"newRoot": bson.M{
					"$arrayElemAt": bson.A{
						"$catalog_info",
						0,
					},
				},
			},
		},
	}

	ctx := context.TODO()

	catalogsCursor, err := gi.DB.Collection(model.GroupColl).Aggregate(ctx, mongo.Pipeline{
		matchStage,
		unwindStage,
		lookupStage,
		skipStage,
		limitStage,
		replaceRootStage,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query for catalog with group id:%s", id.Hex())
	}
	var catalogs []schema.GetCatalogByGroupIDResp
	if err := catalogsCursor.All(ctx, &catalogs); err != nil {
		return nil, errors.Wrap(err, "error decoding Catalogs")
	}

	return catalogs, nil
}

// GetGroups returns a list of Group
func (gi *GroupImpl) GetGroups(opts *schema.GetGroupsOpts) ([]schema.GroupResp, error) {
	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"status.value": opts.Status,
		},
	}}
	limitStage := bson.D{
		{Key: "$limit", Value: gi.App.Config.PageSize},
	}
	skipStage := bson.D{
		{Key: "$skip", Value: gi.App.Config.PageSize * opts.Page},
	}

	lookupStage := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         model.CatalogColl,
			"localField":   "catalog_ids",
			"foreignField": "_id",
			"as":           "catalog_info",
		},
	}}

	setStage := bson.D{{
		Key: "$set", Value: bson.M{
			"minimum": bson.M{
				"$min": "$catalog_info.retail_price.value",
			},
			"maximum": bson.M{
				"$max": "$catalog_info.retail_price.value",
			},
		},
	}}

	ctx := context.TODO()
	var pipeline mongo.Pipeline
	if opts.Status != "all" {
		pipeline = append(pipeline, matchStage)
	}
	pipeline = append(pipeline, mongo.Pipeline{skipStage, limitStage, lookupStage, setStage}...)

	cur, err := gi.DB.Collection(model.GroupColl).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var groupsResp []schema.GroupResp
	if err := cur.All(ctx, &groupsResp); err != nil {
		return nil, err
	}

	return groupsResp, nil
}

//GetGroupsByCatalogID returns List of groups containing that catalog ID.
func (gi *GroupImpl) GetGroupsByCatalogID(opts *schema.GetGroupsByCatalogIDOpts) ([]schema.GetGroupsByCatalogIDResp, error) {
	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"status.value": model.Publish,
			"catalog_ids":  bson.M{"$in": []primitive.ObjectID{opts.ID}},
		},
	}}
	fmt.Println(opts.ID)
	limitStage := bson.D{
		{Key: "$limit", Value: gi.App.Config.PageSize},
	}
	skipStage := bson.D{
		{Key: "$skip", Value: gi.App.Config.PageSize * opts.Page},
	}
	lookupStage := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         model.CatalogColl,
			"localField":   "catalog_ids",
			"foreignField": "_id",
			"as":           "catalog_info",
		},
	}}
	projectStage := bson.D{{
		Key: "$project", Value: bson.M{
			"_id":                         1,
			"basis":                       1,
			"catalog_info._id":            1,
			"catalog_info.name":           1,
			"catalog_info.retail_price":   1,
			"catalog_info.featured_image": 1,
		},
	}}

	ctx := context.TODO()

	cur, err := gi.DB.Collection(model.GroupColl).Aggregate(ctx, mongo.Pipeline{
		matchStage,
		skipStage,
		limitStage,
		lookupStage,
		projectStage,
	})
	if err != nil {
		return nil, err
	}
	var groupsResp []schema.GetGroupsByCatalogIDResp
	if err := cur.All(ctx, &groupsResp); err != nil {
		return nil, err
	}

	return groupsResp, nil

}

//KeeperGetGroupsByCatalogID returns List of groups containing that catalog ID.
func (gi *GroupImpl) KeeperGetGroupsByCatalogID(opts *schema.KeeperGetGroupsByCatalogIDOpts) ([]schema.GroupResp, error) {

	var pipeline mongo.Pipeline

	if opts.Status != "" {
		matchStage := bson.D{{
			Key: "$match", Value: bson.M{
				"status.value": opts.Status,
				"catalog_ids":  bson.M{"$in": []primitive.ObjectID{opts.ID}},
			},
		}}
		pipeline = append(pipeline, matchStage)

	} else {
		matchStage := bson.D{{
			Key: "$match", Value: bson.M{
				"catalog_ids": bson.M{"$in": []primitive.ObjectID{opts.ID}},
			},
		}}
		pipeline = append(pipeline, matchStage)
	}
	limitStage := bson.D{
		{Key: "$limit", Value: gi.App.Config.PageSize},
	}
	skipStage := bson.D{
		{Key: "$skip", Value: gi.App.Config.PageSize * opts.Page},
	}

	lookupStage := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         model.CatalogColl,
			"localField":   "catalog_ids",
			"foreignField": "_id",
			"as":           "catalog_info",
		},
	}}

	setStage := bson.D{{
		Key: "$set", Value: bson.M{
			"minimum": bson.M{
				"$min": "$catalog_info.retail_price.value",
			},
			"maximum": bson.M{
				"$max": "$catalog_info.retail_price.value",
			},
		},
	}}

	ctx := context.TODO()
	pipeline = append(pipeline, mongo.Pipeline{skipStage, limitStage, lookupStage, setStage}...)

	cur, err := gi.DB.Collection(model.GroupColl).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var groupsResp []schema.GroupResp
	if err := cur.All(ctx, &groupsResp); err != nil {
		return nil, err
	}

	return groupsResp, nil

}

//AddCatalogsInTheGroup adds catalogs to the group with the given group id
func (gi *GroupImpl) AddCatalogsInTheGroup(opts *schema.AddCatalogsInTheGroupOpts) (bool, []error) {

	ctx := context.TODO()

	var errorResp []error
	catalogs, err := gi.App.KeeperCatalog.GetCatalogByIDs(ctx, opts.CatalogIDs)
	if err != nil {
		return false, []error{errors.Wrap(err, "unable to fetch catalogs by ids")}
	}
	catalogMap := make(map[primitive.ObjectID]schema.GetCatalogResp)
	for i := 0; i < len(catalogs); i++ {
		catalogMap[catalogs[i].ID] = catalogs[i]
	}

	if len(catalogs) != len(opts.CatalogIDs) {
		for i := 0; i < len(opts.CatalogIDs); i++ {
			_, ok := catalogMap[opts.CatalogIDs[i]]
			if !ok {
				errorResp = append(errorResp, errors.Errorf("catalog with id: %s not found", opts.CatalogIDs[i].Hex()))
			}
		}
		return false, errorResp
	}

	query := bson.M{"$addToSet": bson.M{"catalog_ids": bson.M{"$each": opts.CatalogIDs}}}
	res, err := gi.DB.Collection(model.GroupColl).UpdateOne(ctx, bson.M{"_id": opts.ID}, query)
	if err != nil {
		return false, []error{errors.Wrap(err, "unable to add catalogs to group")}
	}
	if res.MatchedCount == 0 {
		return false, []error{errors.Errorf("group with id:%s not found", opts.ID)}
	}
	if res.ModifiedCount == 0 {
		return false, []error{errors.New("unable to add catalogs to group")}
	}

	return true, nil
}

//RemoveCatalogsFromTheGroup adds catalogs to teh group with the given group id
func (gi *GroupImpl) RemoveCatalogsFromTheGroup(opts *schema.AddCatalogsInTheGroupOpts) (bool, []error) {

	query := bson.M{
		"$pull": bson.M{
			"catalog_ids": bson.M{
				"$in": opts.CatalogIDs,
			},
		},
	}
	res, err := gi.DB.Collection(model.GroupColl).UpdateOne(context.TODO(), bson.M{"_id": opts.ID}, query)
	if err != nil {
		return false, []error{errors.Wrap(err, "unable to remove catalogs from the group")}
	}
	if res.MatchedCount == 0 {
		return false, []error{errors.Errorf("group with id:%s not found", opts.ID)}
	}
	if res.ModifiedCount == 0 {
		return false, []error{errors.New("unable to remove catalogs from the group")}
	}

	return true, nil
}

//UpdateGroupStatus updates the status of the group with the given ID
func (gi *GroupImpl) UpdateGroupStatus(opts *schema.UpdateGroupStatusOpts) error {

	ctx := context.TODO()
	filter := bson.M{
		"_id": opts.ID,
	}

	var group model.Group
	err := gi.DB.Collection(model.GroupColl).FindOne(ctx, filter).Decode(&group)
	if err != nil {
		return errors.Wrapf(err, "failed to find group with id: %s", opts.ID)
	}
	cStatus := group.Status.Value
	if cStatus == model.Archive {
		return errors.Errorf("cannot change status from archive to %s", opts.Status)
	}
	updateStatus := model.GroupStatus{
		Value:     opts.Status,
		CreatedAt: time.Now().UTC(),
	}

	updateQuery := bson.M{
		"$set": bson.M{
			"status": updateStatus,
		},
		"$push": bson.M{
			"status_history": updateStatus,
		},
	}
	if _, err := gi.DB.Collection(model.GroupColl).UpdateOne(ctx, filter, updateQuery); err != nil {
		return errors.Wrap(err, "fail to update Status")
	}
	return nil
}

//GetGroupsByCatalogName gets the group info and catalogs, based on catalog name
func (gi *GroupImpl) GetGroupsByCatalogName(name string, page int) ([]schema.GetGroupsByCatalogNameResp, error) {

	fmt.Println(name)
	fmt.Println(page)

	matchStage := bson.D{{
		Key: "$match", Value: bson.M{
			"lname": bson.M{
				"$regex": name,
			},
		},
	}}

	lookupStage := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         "group",
			"localField":   "_id",
			"foreignField": "catalog_ids",
			"as":           "group_info",
		},
	}}

	projectStage := bson.D{{
		Key: "$project", Value: bson.M{
			"group_info": 1,
			"_id":        0,
		},
	}}
	unwindStage := bson.D{{
		Key: "$unwind", Value: bson.M{
			"path": "$group_info",
		},
	}}

	groupStage := bson.D{{
		Key: "$group", Value: bson.M{
			"_id": "$group_info._id",
			"group_info": bson.M{
				"$first": "$group_info",
			},
		},
	}}
	limitStage := bson.D{
		{Key: "$limit", Value: gi.App.Config.PageSize},
	}
	skipStage := bson.D{
		{Key: "$skip", Value: gi.App.Config.PageSize * page},
	}
	lookupStage2 := bson.D{{
		Key: "$lookup", Value: bson.M{
			"from":         "catalog",
			"localField":   "group_info.catalog_ids",
			"foreignField": "_id",
			"as":           "catalog_info",
		},
	}}

	setStage := bson.D{{
		Key: "$set", Value: bson.M{
			"minimum": bson.M{
				"$min": "$catalog_info.retail_price.value",
			},
			"maximum": bson.M{
				"$max": "$catalog_info.retail_price.value",
			},
		},
	}}

	projectStage2 := bson.D{{
		Key: "$project", Value: bson.M{
			"group_status": "$group_info.status",
			"minimum":      1,
			"maximum":      1,
			"catalog_info": bson.M{
				"$slice": bson.A{"$catalog_info", 3},
			},
			"count": bson.M{
				"$size": "$catalog_info",
			}},
	}}

	ctx := context.TODO()
	var pipeline mongo.Pipeline

	pipeline = append(pipeline, mongo.Pipeline{matchStage, lookupStage, projectStage, unwindStage, groupStage, skipStage, limitStage, lookupStage2, setStage, projectStage2}...)

	cur, err := gi.DB.Collection(model.CatalogColl).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var groupsResp []schema.GetGroupsByCatalogNameResp
	if err := cur.All(ctx, &groupsResp); err != nil {
		return nil, err
	}

	return groupsResp, nil
}

//EditGroup adds catalogs to the group with the given group id
func (gi *GroupImpl) EditGroup(opts *schema.EditGroupOpts) (*schema.EditGroupResp, []error) {

	ctx := context.TODO()
	var errorResp []error
	group := model.Group{}
	catalogs, err := gi.App.KeeperCatalog.GetCatalogByIDs(ctx, opts.CatalogIDs)
	if err != nil {
		return nil, []error{errors.Wrap(err, "unable to fetch catalogs by ids")}
	}
	catalogMap := make(map[primitive.ObjectID]schema.GetCatalogResp)
	for i := 0; i < len(catalogs); i++ {
		catalogMap[catalogs[i].ID] = catalogs[i]
	}

	if len(catalogs) != len(opts.CatalogIDs) {
		for i := 0; i < len(opts.CatalogIDs); i++ {
			_, ok := catalogMap[opts.CatalogIDs[i]]
			if !ok {
				errorResp = append(errorResp, errors.Errorf("catalog with id: %s not found", opts.CatalogIDs[i].Hex()))
			}
		}
		return nil, errorResp
	}

	if opts.Basis != "" {
		group.Basis = opts.Basis
	}
	if len(opts.CatalogIDs) != 0 {
		group.CatalogIDs = opts.CatalogIDs
	}
	updateQuery := bson.M{
		"$set": group,
	}
	qOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = gi.DB.Collection(model.GroupColl).FindOneAndUpdate(ctx, bson.M{"_id": opts.ID}, updateQuery, qOpts).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments || err == mongo.ErrNilDocument {
			return nil, []error{errors.Errorf("catalog with id:%s not found", opts.ID.Hex())}
		}
		return nil, []error{errors.Wrap(err, "unable to replace catalogs to group")}
	}

	return &schema.EditGroupResp{
		ID:         group.ID,
		Basis:      group.Basis,
		Status:     *group.Status,
		CatalogIDs: group.CatalogIDs,
	}, nil
}
