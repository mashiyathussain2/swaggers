package api

// InitRoutes initializes all the endpoints
func (a *API) InitRoutes() {
	a.Router.APIRoot.Handle("/me", a.requestWithAuthHandler(a.me)).Methods("GET")
	a.Router.APIRoot.Handle("/user/forgot-password", a.requestHandler(a.forgotPassword)).Methods("POST")
	a.Router.APIRoot.Handle("/user/reset-password", a.requestHandler(a.resetPassword)).Methods("POST")
	a.Router.APIRoot.Handle("/user/verify-email", a.requestHandler(a.verifyEmail)).Methods("POST")
	a.Router.APIRoot.Handle("/user/verify-email/resend", a.requestHandler(a.resendEmailVerificationCode)).Methods("POST")

	a.Router.APIRoot.Handle("/customer/social/login", a.requestHandler(a.loginViaSocial)).Methods("POST")
	a.Router.APIRoot.Handle("/customer/email/signup", a.requestHandler(a.signUpViaEmail)).Methods("POST")
	a.Router.APIRoot.Handle("/customer/email/login", a.requestHandler(a.loginViaEmail)).Methods("POST")
	a.Router.APIRoot.Handle("/customer/otp/generate", a.requestHandler(a.loginViaMobileOTP)).Methods("POST")
	a.Router.APIRoot.Handle("/customer/otp/confirm", a.requestHandler(a.confirmLoginViaMobileOTP)).Methods("POST")
	a.Router.APIRoot.Handle("/customer", a.requestWithAuthHandler(a.updateCustomerInfo)).Methods("PUT")
}

// InitTestRoutes := intializing all the testing and development endpoints
func (a *API) InitTestRoutes() {
	// a.Router.APIRoot.Handle("/test/add-category", a.requestHandler(a.addSampleCategories)).Methods("GET")
	// a.Router.APIRoot.Handle("/test/add-catalog-and-variants", a.requestHandler(a.imageUpload)).Methods("GET")
}

// InitMediaRoutes initialize media urls
func (a *API) InitMediaRoutes() {
	// fs := http.FileServer(http.Dir("./.media/"))
	// a.Router.Root.PathPrefix("/media/").Handler(http.StripPrefix("/media/", fs))
}
