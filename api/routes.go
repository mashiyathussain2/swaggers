package api

// InitRoutes initializes all the endpoints
func (a *API) InitRoutes() {
	a.Router.Root.Handle("/", a.requestHandler(a.home)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/live/create", a.requestHandler(a.createLiveStream)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/live/{liveID}", a.requestHandler(a.getLiveStreamByID)).Methods("GET")
	a.Router.APIRoot.Handle("/keeper/live", a.requestHandler(a.getLiveStreams)).Methods("GET")

	a.Router.APIRoot.Handle("/keeper/content", a.requestHandler(a.getContent)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/content/get", a.requestHandler(a.geContents)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/content/status", a.requestHandler(a.changeContentStatus)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/content/pebble", a.requestHandler(a.createPebble)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/content/pebble", a.requestHandler(a.editPebble)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/content/pebble/process", a.requestHandler(a.processPebble)).Methods("PUT")
	a.Router.APIRoot.Handle("/keeper/content/{contentID}", a.requestHandler(a.getContentByID)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/content/pebble/{pebbleID}", a.requestHandler(a.deletePebble)).Methods("DELETE")
	a.Router.APIRoot.Handle("/keeper/content/catalog/video", a.requestHandler(a.createVideoCatalogContent)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/content/catalog/image", a.requestHandler(a.createImageCatalogContent)).Methods("POST")
	a.Router.APIRoot.Handle("/keeper/content/catalog", a.requestHandler(a.editCatalogContent)).Methods("POST")
	a.Router.APIRoot.Handle("/image/upload", a.requestHandler(a.uploadImage)).Methods("POST")

	a.Router.APIRoot.Handle("/live/{liveID}/start", a.requestHandler(a.startLiveStream)).Methods("GET")
	a.Router.APIRoot.Handle("/live/{liveID}/stop", a.requestHandler(a.stopLiveStream)).Methods("GET")
	a.Router.APIRoot.Handle("/live/{liveID}/join", a.requestHandler(a.joinLiveStream)).Methods("GET")
	a.Router.APIRoot.Handle("/live/{liveID}/comment", a.requestHandler(a.pushComment)).Methods("POST")

	a.Router.APIRoot.Handle("/content/like", a.requestHandler(a.createLike)).Methods("POST")
	a.Router.APIRoot.Handle("/content/view", a.requestHandler(a.createView)).Methods("POST")
	a.Router.APIRoot.Handle("/content/comment", a.requestHandler(a.createContentComment)).Methods("POST")

	a.Router.APIRoot.Handle("/pebble", a.requestHandler(a.getPebble)).Methods("GET")

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
