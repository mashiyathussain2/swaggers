package middleware

import (
	"go-app/server/auth"
	"net/http"
)

// AuthenticationMiddleware to initialize authentication middleware
type AuthenticationMiddleware struct {
	Session auth.SessionAuth
}

// NewAuthenticationMiddleware creates an instance to authentication middleware
func NewAuthenticationMiddleware(s auth.SessionAuth) *AuthenticationMiddleware {
	authMiddleware := AuthenticationMiddleware{
		Session: s,
	}
	return &authMiddleware
}

// GetMiddlewareHandler get cookie and sets the token to authorization header of request
func (am *AuthenticationMiddleware) GetMiddlewareHandler() func(http.ResponseWriter, *http.Request, http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if r.Header.Get("Authorization") == "" {
			uc, _ := am.Session.Get(r)
			if uc != nil && uc.Token != "" {
				r.Header.Set("Authorization", uc.Token)
			}
		}
		next(rw, r)
	}
}
