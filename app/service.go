package app

import "fmt"

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

	fmt.Println(a.Config)
	a.Live = InitLive(&LiveOpts{
		App:    a,
		DB:     a.MongoDB.Client.Database(a.Config.LiveConfig.DBName),
		Logger: a.Logger,
		IVS:    NewIVSImpl(&IVSOpts{Config: &a.Config.IVSConfig}),
	})
}
