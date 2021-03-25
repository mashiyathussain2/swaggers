package app

// InitConsumer initializes consumers
func InitConsumer(a *App) {
	// ctx := context.Background()
	// a.CatalogListener = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
	// 	Logger: a.Logger,
	// 	Config: &a.Config.CatalogListenerConfig,
	// })
	// go a.CatalogListener.ConsumeAndCommit(ctx, func(m kafka.Message) { fmt.Println(m) })
}

// CloseConsumer close all consumer connections
func CloseConsumer(a *App) {
	// a.CatalogListener.Close()
}
