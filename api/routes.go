package api

// InitRoutes initializes all the endpoints
func (a *API) InitRoutes() {
	a.Router.Root.Handle("/keeper/gl/callback", a.requestHandler(a.keeperLoginCallback)).Methods("GET")

	a.Router.APIRoot.Handle("/me", a.requestWithAuthHandler(a.me)).Methods("GET")

	a.Router.APIRoot.Handle("/keeper/auth/login", a.requestHandler(a.keeperLogin)).Methods("GET")

	a.Router.APIRoot.Handle("/keeper/user/get", a.requestWithSudoHandler(a.getUserInfoByID)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/brand", a.requestWithSudoHandler(a.createbrand)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/brands", a.requestWithSudoHandler(a.getBrands)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/brand/get", a.requestWithSudoHandler(a.getBrandsById)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/brand", a.requestWithSudoHandler(a.editbrand)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/brand/{brandID}/check", a.requestWithSudoHandler(a.checkBrandByID)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/influencer", a.requestWithSudoHandler(a.createInfluencer)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/influencer/get", a.requestWithSudoHandler(a.getInfluencersByID)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/influencer/name/get", a.requestWithSudoHandler(a.getInfluencerByName)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/influencer", a.requestWithSudoHandler(a.editInfluencer)).Methods("PUT")

	a.Router.APIRoot.Handle("/user/forgot-password", a.requestHandler(a.forgotPassword)).Methods("POST")
	a.Router.APIRoot.Handle("/user/reset-password", a.requestHandler(a.resetPassword)).Methods("POST")
	a.Router.APIRoot.Handle("/user/verify-email", a.requestHandler(a.verifyEmail)).Methods("POST")
	a.Router.APIRoot.Handle("/user/verify-email/resend", a.requestHandler(a.resendEmailVerificationCode)).Methods("POST")

	// LOGIN AND SIGNUP APIS
	a.Router.APIRoot.Handle("/customer/social/login", a.requestHandler(a.loginViaSocial)).Methods("POST")
	a.Router.APIRoot.Handle("/customer/email/signup", a.requestHandler(a.signUpViaEmail)).Methods("POST")
	a.Router.APIRoot.Handle("/customer/email/login", a.requestHandler(a.loginViaEmail)).Methods("POST")
	a.Router.APIRoot.Handle("/customer/otp/generate", a.requestHandler(a.loginViaMobileOTP)).Methods("POST")
	a.Router.APIRoot.Handle("/customer/otp/confirm", a.requestHandler(a.confirmLoginViaMobileOTP)).Methods("POST")

	// CUSTOMER APIS
	a.Router.APIRoot.Handle("/customer", a.requestWithAuthHandler(a.updateCustomerInfo)).Methods("PUT")
	a.Router.APIRoot.Handle("/customer/address", a.requestWithAuthHandler(a.addAddress)).Methods("PUT")
	a.Router.APIRoot.Handle("/customer/{userID}/address", a.requestWithAuthHandler(a.getAddress)).Methods("GET")

	// TODO: Shall i remove this api??
	a.Router.APIRoot.Handle("/brand/{brandID}", a.requestWithAuthHandler(a.getBrandByID)).Methods("GET")

	a.Router.APIRoot.Handle("/app/cart", a.requestWithAuthHandler(a.addToCart)).Methods("POST")
	a.Router.APIRoot.Handle("/app/cart/item", a.requestWithAuthHandler(a.updateItemQty)).Methods("PUT")
	// a.Router.APIRoot.Handle("/app/cart/{userID}", a.requestHandler(a.createCart)).Methods("POST")
	a.Router.APIRoot.Handle("/app/cart/{userID}", a.requestWithAuthHandler(a.getCartInfo)).Methods("GET")
	a.Router.APIRoot.Handle("/app/cart/address", a.requestWithAuthHandler(a.setCartAddress)).Methods("POST")
	a.Router.APIRoot.Handle("/app/cart/{cartID}/checkout", a.requestWithAuthHandler(a.checkoutCart)).Methods("GET")

	a.Router.APIRoot.Handle("/app/customer/{customerID}", a.requestWithAuthHandler(a.getCustomerInfo)).Methods("GET")
	a.Router.APIRoot.Handle("/app/customer/influencer/follow", a.requestWithAuthHandler(a.followInfluencer)).Methods("POST")
	a.Router.APIRoot.Handle("/app/customer/brand/follow", a.requestWithAuthHandler(a.followBrand)).Methods("POST")

	a.Router.APIRoot.Handle("/app/brand/basic", a.requestWithAuthHandler(a.getBrandsBasic)).Methods("POST")
	a.Router.APIRoot.Handle("/app/brand/{brandID}", a.requestWithAuthHandler(a.getBrandInfo)).Methods("GET")

	a.Router.APIRoot.Handle("/app/influencer/basic", a.requestWithAuthHandler(a.getInfluencersBasic)).Methods("POST")
	a.Router.APIRoot.Handle("/app/influencer/{influencerID}", a.requestWithAuthHandler(a.getInfluencerInfo)).Methods("GET")
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
