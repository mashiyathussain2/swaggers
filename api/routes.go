package api

// InitRoutes initializes all the endpoints
func (a *API) InitRoutes() {
	a.Router.Root.Handle("/", a.requestHandler(a.home)).Methods("GET")

	a.Router.APIRoot.Handle("/keeper/content/get", a.requestWithSudoHandler(a.geContents)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/content/status", a.requestWithSudoHandler(a.changeContentStatus)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/content/pebble", a.requestWithSudoHandler(a.createPebble)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/content/pebble", a.requestWithSudoHandler(a.editPebble)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/content/pebble/process", a.requestWithInternalHandler(a.processPebble)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/content/catalog", a.requestWithSudoHandler(a.editCatalogContent)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/content/{contentID}", a.requestWithSudoHandler(a.getContentByID)).Methods("POST")
	a.Router.APIRoot.Handle("/image/upload", a.requestHandler(a.uploadImage)).Methods("POST")

	a.Router.APIRoot.Handle("/keeper/content", a.requestWithInternalHandler(a.getContent)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/content/{contentID}", a.requestWithInternalHandler(a.deleteContent)).Methods("DELETE")
	a.Router.APIRoot.Handle("/keeper/content/catalog/video", a.requestWithInternalHandler(a.createVideoCatalogContent)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/content/catalog/image", a.requestWithInternalHandler(a.createImageCatalogContent)).Methods("POST")

	a.Router.APIRoot.Handle("/keeper/content/review/video", a.requestWithInternalHandler(a.createVideoReviewContent)).Methods("POST")

	a.Router.APIRoot.Handle("/keeper/live", a.requestWithSudoHandler(a.getLiveStreams)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/live/create", a.requestWithSudoHandler(a.createLiveStream)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/live/{liveID}", a.requestWithSudoHandler(a.getLiveStreamByID)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/live/{liveID}/catalog", a.requestWithSudoHandler(a.pushCatalog)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/live/{liveID}/order", a.requestWithSudoHandler(a.pushOrder)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/live/{liveID}/start", a.requestWithSudoHandler(a.startLiveStream)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/live/{liveID}/stop", a.requestWithSudoHandler(a.stopLiveStream)).Methods("GET")

	a.Router.APIRoot.Handle("/live", a.requestHandler(a.getAppLiveStreams)).Methods("GET")
	a.Router.APIRoot.Handle("/live/{liveID}", a.requestHandler(a.getAppLiveStreamByID)).Methods("GET")
	a.Router.APIRoot.Handle("/live/{liveID}/join", a.requestHandler(a.joinLiveStream)).Methods("GET")
	a.Router.APIRoot.Handle("/live/{liveID}/joined", a.requestHandler(a.joinedLiveStream)).Methods("POST")
	a.Router.APIRoot.Handle("/live/{liveID}/comment", a.requestHandler(a.pushComment)).Methods("POST")

	a.Router.APIRoot.Handle("/content/like", a.requestWithAuthHandler(a.createLike)).Methods("POST")
	a.Router.APIRoot.Handle("/content/view", a.requestWithAuthHandler(a.createView)).Methods("POST")
	a.Router.APIRoot.Handle("/content/comment", a.requestWithAuthHandler(a.createContentComment)).Methods("POST")

	a.Router.APIRoot.Handle("/pebble", a.requestHandler(a.getPebble)).Methods("GET")
	a.Router.APIRoot.Handle("/pebble/id", a.requestHandler(a.getPebbleByID)).Methods("GET")
	a.Router.APIRoot.Handle("/pebble/brand", a.requestHandler(a.getPebblesByBrandID)).Methods("GET")
	a.Router.APIRoot.Handle("/pebble/influencer", a.requestHandler(a.getPebblesByInfluencerID)).Methods("GET")

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
