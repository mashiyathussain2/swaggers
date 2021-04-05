package api

import (
	"fmt"
	"go-app/schema"
	"go-app/server/handler"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vasupal1996/goerror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) createCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateCollectionOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, errs := a.App.Collection.CreateCollection(&s)
	if errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}
func (a *API) deleteCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["collectionID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["collectionID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	err = a.App.Collection.DeleteCollection(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
func (a *API) addSubCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddSubCollectionOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, errs := a.App.Collection.AddSubCollection(&s)
	if errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
}
func (a *API) deleteSubCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	collID, err := primitive.ObjectIDFromHex(mux.Vars(r)["collectionID"])

	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["collectionID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	subCollID, err := primitive.ObjectIDFromHex(mux.Vars(r)["subCollectionID"])

	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid subcollection id:%s in url", mux.Vars(r)["subCollectionID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	err = a.App.Collection.DeleteSubCollection(collID, subCollID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
func (a *API) editCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditCollectionOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Collection.EditCollection(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}
func (a *API) updateSubCollectionImage(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateSubCollectionImageOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.Collection.UpdateSubCollectionImage(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
func (a *API) addCatalogsToSubCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateCatalogsInSubCollectionOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	errs := a.App.Collection.AddCatalogsToSubCollection(&s)
	if errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
func (a *API) removeCatalogsFromSubCollection(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateCatalogsInSubCollectionOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	errs := a.App.Collection.RemoveCatalogsFromSubCollection(&s)
	if errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) getCollections(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	page := GetPageValue(r)

	resp, err := a.App.Collection.GetCollections(page)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)

}

func (a *API) getActiveCollections(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	resp, err := a.App.Elasticsearch.GetActiveCollections()
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)

}

func (a *API) updateCollectionStatus(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateCollectionStatus
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.Collection.UpdateCollectionStatus(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
