package api

import (
	"go-app/schema"
	"go-app/server/handler"
	"net/http"

	"github.com/pkg/errors"
	"github.com/vasupal1996/goerror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) createCatalogGroup(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var req schema.CreateCatalogGroupOpts
	if err := a.DecodeJSONBody(r, &req); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&req); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	_, err := a.App.Group.CreateCatalogGroup(&req)
	if err != nil {
		requestCTX.SetErrs(err, http.StatusBadRequest)
		return
	}

	requestCTX.SetAppResponse(true, http.StatusCreated)
}
func (a *API) getCatalogsByGroupID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	page := GetPageValue(r)
	id, err := GetObjectID(r)
	if err != nil {
		requestCTX.SetErr(errors.Wrapf(err, "unable to parse id: %s", r.URL.Query().Get("id")), http.StatusBadRequest)
		return
	}
	if id == primitive.NilObjectID {
		requestCTX.SetErr(goerror.New("id cannot be empty", &goerror.BadRequest), http.StatusBadRequest)
		return
	}

	res, err := a.App.Group.GetCatalogsByGroupID(id, page)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getGroups(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	page := GetPageValue(r)
	status := GetStatusValue(r)

	s := schema.GetGroupsOpts{
		Page:   page,
		Status: status,
	}

	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Group.GetGroups(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}
func (a *API) getGroupsByCatalogID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	page := GetPageValue(r)
	id, err := GetObjectID(r)
	if err != nil {
		requestCTX.SetErr(errors.Wrapf(err, "unable to parse id: %s", r.URL.Query().Get("id")), http.StatusBadRequest)
		return
	}
	if id == primitive.NilObjectID {
		requestCTX.SetErr(goerror.New("id cannot be empty", &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	s := schema.GetGroupsByCatalogIDOpts{
		ID:   id,
		Page: page,
	}
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Group.GetGroupsByCatalogID(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) keeperGetGroupsByCatalogID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	id, err := GetObjectID(r)
	if err != nil {
		requestCTX.SetErr(errors.Wrapf(err, "unable to parse id: %s", r.URL.Query().Get("id")), http.StatusBadRequest)
		return
	}
	if id == primitive.NilObjectID {
		requestCTX.SetErr(goerror.New("id cannot be empty", &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	page := GetPageValue(r)
	status := GetStatusValue(r)
	s := schema.KeeperGetGroupsByCatalogIDOpts{
		ID:     id,
		Page:   page,
		Status: status,
	}

	res, err := a.App.Group.KeeperGetGroupsByCatalogID(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) addCatalogsInTheGroup(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddCatalogsInTheGroupOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Group.AddCatalogsInTheGroup(&s)
	if err != nil {
		requestCTX.SetErrs(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) editGroup(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditGroupOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Group.EditGroup(&s)
	if err != nil {
		requestCTX.SetErrs(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)

}

func (a *API) getGroupsByCatalogName(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {

	page := GetPageValue(r)
	name := r.URL.Query().Get("name")
	res, err := a.App.Group.GetGroupsByCatalogName(name, page)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)

}
