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

func (a *API) addToWishlist(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.AddToWishlistOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.Wishlist.AddToWishlist(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}

func (a *API) getWishlist(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["userID"])
	if err != nil {
		requestCTX.SetErr(goerror.New(fmt.Sprintf("invalid id:%s in url", mux.Vars(r)["userID"]), &goerror.BadRequest), http.StatusBadRequest)
		return
	}
	res, err := a.App.Wishlist.GetWishlist(id)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(res, http.StatusOK)
}

func (a *API) removeFromWishlist(requestCTX *handler.RequestContext, w http.ResponseWriter, r *http.Request) {
	var s schema.RemoveFromWishlistOpts
	if err := a.DecodeJSONBody(r, &s); err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	if errs := a.Validator.Validate(&s); errs != nil {
		requestCTX.SetErrs(errs, http.StatusBadRequest)
		return
	}
	err := a.App.Wishlist.RemoveFromWishlist(&s)
	if err != nil {
		requestCTX.SetErr(err, http.StatusBadRequest)
		return
	}
	requestCTX.SetAppResponse(true, http.StatusOK)
}
