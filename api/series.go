package api

import (
	"go-app/schema"
	"go-app/server/auth"
	"go-app/server/handler"
	"net/http"

	"github.com/pasztorpisti/qs"
)

func (a *API) createPebbleSeries(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateSeriesOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.Series.CreateSeries(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusCreated)
	return
}

func (a *API) editPebbleSeries(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateSeriesOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.Series.UpdateSeries(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusCreated)
	return
}

func (a *API) getPebbleSeries(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetPebbleSeriesFilter
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.UserID = requestCTX.UserClaim.(*auth.UserClaim).ID
	}
	res, err := a.App.Elasticsearch.GetPebbleSeries(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) getPebbleSeriesByIDs(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetSeriesByIDs
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.UserID = requestCTX.UserClaim.(*auth.UserClaim).ID
	}
	res, err := a.App.Elasticsearch.GetSeriesByIDs(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) keeperGetPebbleSeries(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetContentFilter
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Series.GetContentForPebbleSeries(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) keeperGetSeries(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetSeriesKeeperFilter
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Series.KeeperGetSeries(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) keeperGetPebbleSeriesByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.KeeperGetSeriesByID
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	res, err := a.App.Series.KeeperGetSeriesByID(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) keeperGetSeriesBasic(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.KeeperGetSeriesBasic
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Series.KeeperGetSeriesBasic(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}
