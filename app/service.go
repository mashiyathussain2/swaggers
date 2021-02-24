package app

// InitService this initializes all the busines logic services
func InitService(a *App) {
	a.Media = InitContent(&MediaImplOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.ContentConfig.DBName),
		Logger: a.Logger,
		S3: InitS3(&S3Opts{
			Config: &a.Config.S3Config,
		}),
	})

	a.Content = InitPebble(&PebbleOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.PebbleConfig.DBName),
		Logger: a.Logger,
	})
}
