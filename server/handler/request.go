package handler

import (
	"encoding/json"
	"fmt"
	"go-app/server/auth"
	"go-app/server/middleware"
	"net/http"

	errors "github.com/vasupal1996/goerror"
)

// Request represents a request from client
type Request struct {
	HandlerFunc func(*RequestContext, http.ResponseWriter, *http.Request)
	AuthFunc    auth.TokenAuth
	IsLoggedIn  bool
	IsSudoUser  bool
}

// HandleRequest := handles incoming requests from client
func (rh *Request) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestCTX := &RequestContext{}
	requestCTX.RequestID = middleware.RequestIDFromContext(r.Context())
	requestCTX.Path = r.URL.Path

	authToken := r.Header.Get("Authorization")
	if authToken != "" {
		claim, err := rh.AuthFunc.VerifyToken(authToken)
		if err != nil {
			requestCTX.SetErr(errors.New("failed to verify token", &errors.PermissionDenied), http.StatusUnauthorized)
			goto SKIP_REQUEST
		} else {
			requestCTX.UserClaim = claim.(*auth.UserClaim)
		}
	}

	if rh.IsLoggedIn {
		if requestCTX.UserClaim == nil {
			requestCTX.SetErr(errors.New("auth token required", &errors.PermissionDenied), http.StatusUnauthorized)
			goto SKIP_REQUEST
		} else {
			if rh.IsSudoUser {
				if !requestCTX.UserClaim.IsSudo() {
					requestCTX.SetErr(errors.New("permission denied: required keeper user type", &errors.PermissionDenied), http.StatusForbidden)
					goto SKIP_REQUEST
				}
				fmt.Println("Sudo user")
				cookie, err := r.Cookie("session")
				if err != nil {
					requestCTX.SetErr(errors.Wrap(err, "failed to get session id", &errors.BadRequest), http.StatusBadGateway)
					goto SKIP_REQUEST
				}
				fmt.Println(cookie)
				if cookie.Value != "" {
					err := rh.AuthFunc.AuthorizeKeeperRequest(r.Method, r.Host, r.RequestURI, cookie.Value)
					fmt.Println(err)
					if err != nil {
						requestCTX.SetErr(err, http.StatusUnauthorized)
						goto SKIP_REQUEST
					}
				}
			} else {
				if requestCTX.UserClaim.IsSudo() {
					requestCTX.SetErr(errors.New("permission denied: required customer type", &errors.PermissionDenied), http.StatusForbidden)
					goto SKIP_REQUEST
				}
			}
		}
	} else {
		if rh.IsSudoUser {
			if requestCTX.UserClaim == nil {
				requestCTX.SetErr(errors.New("auth token required", &errors.PermissionDenied), http.StatusUnauthorized)
				goto SKIP_REQUEST
			}
			if !requestCTX.UserClaim.IsInternal() {
				requestCTX.SetErr(errors.New("permission denied: must be internal-user", &errors.PermissionDenied), http.StatusForbidden)
				goto SKIP_REQUEST
			}
		}
	}

SKIP_REQUEST:

	w.Header().Set(auth.HeaderRequestID, requestCTX.RequestID)
	if requestCTX.Err == nil {
		rh.HandlerFunc(requestCTX, w, r)
	}

	if requestCTX.ResponseCode != 0 && requestCTX.ResponseType != RedirectResp {
		w.WriteHeader(requestCTX.ResponseCode)
	}

	switch t := requestCTX.ResponseType; t {
	case HTMLResp:
		w.Header().Set("Content-Type", "text/html")
		res := requestCTX.Response.GetRaw()
		w.Write(res.([]byte))
	case JSONResp:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(requestCTX.Response)
	case ErrorResp:
		w.Header().Set("Content-Type", "application/json")
		requestCTX.Err.RequestID = &requestCTX.RequestID
		json.NewEncoder(w).Encode(&requestCTX.Err)
	case RedirectResp:
		payload := requestCTX.Response.(*AppResponse).Payload
		http.Redirect(w, r, payload.(string), requestCTX.ResponseCode)
	}
}
