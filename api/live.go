package api

import (
	"fmt"
	"go-app/schema"
	"go-app/server/handler"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pasztorpisti/qs"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *API) getLiveStreams(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetLiveStreamsFilter
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	fmt.Println(r.URL.Query().Encode())
	res, err := a.App.Live.GetLiveStreams(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) getLiveStreamByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["liveID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid live id: %s in url", mux.Vars(r)["liveID"]), http.StatusBadRequest)
		return
	}
	res, err := a.App.Live.GetLiveStreamByID(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) createLiveStream(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateLiveStreamOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Live.CreateLiveStream(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) pushComment(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateLiveCommentOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	a.App.Live.PushComment(&s)
	requestCTX.SetAppResponse(true, http.StatusCreated)
	return
}

func (a *API) startLiveStream(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["liveID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid live id: %s in url", mux.Vars(r)["content"]), http.StatusBadRequest)
		return
	}

	res, err := a.App.Live.StartLiveStream(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) joinLiveStream(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["liveID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid live id: %s in url", mux.Vars(r)["content"]), http.StatusBadRequest)
		return
	}

	res, err := a.App.Live.JoinLiveStream(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) stopLiveStream(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["liveID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid live id: %s in url", mux.Vars(r)["content"]), http.StatusBadRequest)
		return
	}

	if err := a.App.Live.EndLiveStream(id); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	requestCTX.SetAppResponse(true, http.StatusCreated)
	return
}

func (a *API) pushCatalog(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.PushCatalogOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	a.App.Live.PushCatalog(&s)
	requestCTX.SetAppResponse(true, http.StatusCreated)
	return
}

func (a *API) pushOrder(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.PushNewOrderOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	a.App.Live.PushOrder(&s)
	requestCTX.SetAppResponse(true, http.StatusCreated)
	return
}
