package app

// InitService this initializes all the busines logic services
func InitService(a *App) {
	a.Media = InitMedia(&MediaImplOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.MediaConfig.DBName),
		Logger: a.Logger,
	})

	a.Content = InitContent(&ContentOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.ContentConfig.DBName),
		Logger: a.Logger,
	})
}
