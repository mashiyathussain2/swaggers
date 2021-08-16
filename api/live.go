package api

import (
	"go-app/schema"
	"go-app/server/auth"
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
	if s.UserID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
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

func (a *API) joinedLiveStream(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.PushJoinOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.ID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).ID {
		requestCTX.SetErr(errors.New("invalid user"), http.StatusForbidden)
		return
	}
	a.App.Live.PushJoin(&s)
	requestCTX.SetAppResponse(true, http.StatusCreated)
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

func (a *API) getAppLiveStreams(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetAppLiveStreamsFilter
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	res, err := a.App.Live.GetAppLiveStreams(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) getAppLiveStreamByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["liveID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid live id: %s in url", mux.Vars(r)["liveID"]), http.StatusBadRequest)
		return
	}
	res, err := a.App.Live.GetAppLiveStreamByID(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) getAppLiveStreamsByInfluencerID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetAppLiveStreamsFilter
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	influencer_id, err := primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if influencer_id == primitive.NilObjectID {
		requestCTX.SetErr(errors.New("influencer id missing"), http.StatusBadRequest)
		return
	}
	res, err := a.App.Live.GetAppLiveStreamsByInfluencerID(influencer_id, &s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusCreated)
	return
}

func (a *API) createLiveStreamByApp(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateLiveStreamOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	influencer_id, err := primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if influencer_id == primitive.NilObjectID {
		requestCTX.SetErr(errors.New("influencer id missing from profile"), http.StatusBadRequest)
		return
	}
	s.InfluencerIDs = append(s.InfluencerIDs, influencer_id)
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
