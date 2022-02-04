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

func (a *API) createInfluencer(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CreateInfluencerOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.CreateInfluencer(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
	return
}

func (a *API) editInfluencer(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditInfluencerOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.EditInfluencer(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getInfluencersByID(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencersByIDOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.GetInfluencersByID(s.IDs)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getInfluencerByName(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencersByNameOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.GetInfluencerByName(s.Name)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getInfluencersBasic(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencersByIDBasicOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.CustomerID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).CustomerID)
	}
	res, err := a.App.Elasticsearch.GetInfluencersByIDBasic(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getInfluencerInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["influencerID"])
	if err != nil {
		requestCTX.SetErr(errors.Errorf("invalid influencer id:%s in url", mux.Vars(r)["influencerID"]), http.StatusBadRequest)
		return
	}
	var userID primitive.ObjectID
	if requestCTX.UserClaim != nil {
		userID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).CustomerID)
	}
	res, err := a.App.Elasticsearch.GetInfluencerInfoByID(&schema.GetInfluencerInfoByIDOpts{ID: id, CustomerID: userID})
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) claimInfluencerRequest(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.InfluencerAccountRequestOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.UserID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)

	}

	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if err := a.App.Influencer.InfluencerAccountRequest(&s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) checkClaimInfluencerRequestStatus(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var userID primitive.ObjectID
	if requestCTX.UserClaim != nil {
		userID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
	}
	res, err := a.App.Influencer.GetInfluencerAccountRequestStatus(userID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getInfluencerClaimRequests(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	res, err := a.App.Influencer.GetInfluencerAccountRequest()
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) updateClaimInfluencerRequestStatus(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateInfluencerAccountRequestStatusOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.GranteeID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
	}
	if err := a.App.Influencer.UpdateInfluencerAccountRequestStatus(&s); err != nil {
		a.Logger.Err(err).Msg("failed to update status request")
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) checkInfluencerUsernameExists(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		requestCTX.SetErr(errors.Errorf("username cannot be empty nil"), http.StatusBadRequest)
		return
	}
	err := a.App.Influencer.CheckInfluencerUsernameExists(username, nil)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) getInfluencersBasicByUsername(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencersByUsernameBasicOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.CustomerID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).CustomerID)
	}
	res, err := a.App.Elasticsearch.GetInfluencersByUserameBasic(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getInfluencerInfoByUsername(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	if username == "" {
		requestCTX.SetErr(errors.Errorf("invalid influencer id:%s in url", mux.Vars(r)["influencerID"]), http.StatusBadRequest)
		return
	}
	var userID primitive.ObjectID
	if requestCTX.UserClaim != nil {
		userID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).CustomerID)
	}
	res, err := a.App.Elasticsearch.GetInfluencerInfoByUsername(&schema.GetInfluencerInfoByUsernameOpts{Username: username, CustomerID: userID})
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) editInfluencerApp(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditInfluencerAppOpts

	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}

	if s.ID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID {
		requestCTX.SetErr(errors.New("not authorized"), http.StatusForbidden)
		return
	}

	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.EditInfluencerApp(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) debitRequest(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.CommissionDebitRequest
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.ID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if err := a.App.Influencer.DebitRequest(&s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
func (a *API) updateDebitRequest(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.UpdateCommissionDebitRequest
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.GranteeID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
	}

	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if err := a.App.Influencer.UpdateDebitRequest(&s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) getDebitRequest(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	res, err := a.App.Influencer.GetActiveDebitRequest()
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getInfluencerDashboard(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencerDashboardOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.GetInfluencerDashboard(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getInfluencerLedger(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetInfluencerLedgerOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	s.ID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.GetInfluencerLedger(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getInfluencerPayoutInfo(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.GetInfluencerPayoutInfo(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) getCommissionAndRevenue(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.GetCommissionAndRevenueOpts
	if err := qs.Unmarshal(&s, r.URL.Query().Encode()); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	s.ID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.GetCommissionAndRevenue(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) editInfluencerAppV2(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.EditInfluencerAppV2Opts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	var err error
	s.ID, err = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID)
	if err != nil {
		requestCTX.SetErr(errors.Wrapf(err, "influencer id invalid"), http.StatusBadRequest)
		return
	}
	if s.ID.Hex() != requestCTX.UserClaim.(*auth.UserClaim).InfluencerInfo.ID {
		requestCTX.SetErr(errors.New("not authorized"), http.StatusForbidden)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	res, err := a.App.Influencer.EditInfluencerAppV2(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) claimInfluencerRequestV2(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.InfluencerAccountRequestV2Opts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if requestCTX.UserClaim != nil {
		s.UserID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).ID)
		s.CustomerID, _ = primitive.ObjectIDFromHex(requestCTX.UserClaim.(*auth.UserClaim).CustomerID)
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	if s.Email != "" {
		err := a.App.User.UpdateUserEmail(&schema.UpdateUserEmailOpts{
			ID:    s.UserID,
			Email: s.Email,
		})
		if err != nil {
			requestCTX.SetErr(err, http.StatusBadRequest)
			return
		}
	}
	if s.Phone != nil {
		err := a.App.User.UpdateUserPhoneNo(&schema.UpdateUserPhoneNoOpts{
			ID:      s.UserID,
			PhoneNo: s.Phone,
		})
		if err != nil {
			requestCTX.SetErr(err, http.StatusBadRequest)
			return
		}
	}
	if err := a.App.Influencer.InfluencerAccountRequestV2(&s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
