package api

import (
	"go-app/schema"
	"go-app/server/handler"
	"net/http"
)

func (a *API) createCatalog(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateCatalogOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	resp, err := a.App.KeeperCatalog.CreateCatalog(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusInternalServerError)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) editCatalog(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditCatalogOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	resp, err := a.App.KeeperCatalog.EditCatalog(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusInternalServerError)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getBasicCatalogInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetBasicCatalogFilter
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	resp, err := a.App.KeeperCatalog.GetBasicCatalogInfo(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusInternalServerError)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getCatalogFilter(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	resp, err := a.App.KeeperCatalog.GetCatalogFilter()
	if err != nil {
		requestCTX.SetErr(err, http.StatusInternalServerError)
		return
	}
	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getCatalogBySlug(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	requestCTX.SetAppResponse("ok", http.StatusOK)
}
