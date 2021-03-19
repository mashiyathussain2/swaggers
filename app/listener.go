package app

// InitConsumer initializes all kafka consumers
func InitConsumer(a *App) {}

// CloseConsumer terminates all consumer connections
func CloseConsumer(a *App) {}

// InitProducer initializes kafka message producers
func InitProducer(a *App) {}

// CloseProducer terminates all producer connections
func CloseProducer(a *App) {}
