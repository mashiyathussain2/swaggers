package api

import (
	"go-app/schema"
	"go-app/server/handler"
	"net/http"

	"github.com/gorilla/mux"
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
	res, err := a.App.KeeperCatalog.CreateCatalog(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
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
	res, err := a.App.KeeperCatalog.EditCatalog(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) addVariants(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddVariantOpts

	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.KeeperCatalog.AddVariant(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) keeperSearchCatalog(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.KeeperSearchCatalogOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.KeeperCatalog.KeeperSearchCatalog(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) deleteVariant(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.DeleteVariantOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.KeeperCatalog.DeleteVariant(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) updateCatalogStatus(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateCatalogStatusOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	resp, err := a.App.KeeperCatalog.UpdateCatalogStatus(&s)
	if err != nil {
		if resp != nil {
			requestCTX.SetCustomResponse(false, resp, err, http.StatusBadRequest)
			return
		}
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) addCatalogContent(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddCatalogContentOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	resp, errs := a.App.KeeperCatalog.AddCatalogContent(&s)
	if errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}

	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) addCatalogContentImage(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddCatalogContentImageOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	errs := a.App.KeeperCatalog.AddCatalogContentImage(&s)
	if errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}

	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) getCatalogsByFilter(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetCatalogsByFilterOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	resp, err := a.App.KeeperCatalog.GetCatalogsByFilter(&s)
	if err != nil {

		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	requestCTX.SetAppResponse(resp, http.StatusOK)
}

func (a *API) getCatalogBySlug(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	slug := mux.Vars(r)["slug"]
	resp, err := a.App.KeeperCatalog.GetCatalogBySlug(slug)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
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
