package app

import (
	"context"
	"go-app/server/kafka"
)

// InitConsumer initializes all kafka consumers
func InitConsumer(a *App) {
	ctx := context.TODO()

	a.CustomerChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.CustomerChangeConfig,
	})
	go a.CustomerChanges.Consume(ctx, a.UserProcessor.ProcessCustomerUpdate)

	a.DiscountChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.DiscountChangeConfig,
	})
	go a.DiscountChanges.Consume(ctx, a.CartProcessor.ProcessDiscountUpdate)

	a.CatalogChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.CatalogChangeConfig,
	})
	go a.CatalogChanges.Consume(ctx, a.CartProcessor.ProcessCatalogUpdate)

	a.InventoryChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.InventoryChangeConfig,
	})
	go a.InventoryChanges.Consume(ctx, a.CartProcessor.ProcessDiscountUpdate)

	a.BrandChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.BrandChangeConfig,
	})
	go a.BrandChanges.ConsumeAndCommit(ctx, a.BrandProcessor.ProcessBrandUpdate)

	a.InfluencerChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.InfluencerChangeConfig,
	})
	go a.InfluencerChanges.ConsumeAndCommit(ctx, a.InfluencerProcessor.ProcessInfluencerUpdate)
}

// CloseConsumer terminates all consumer connections
func CloseConsumer(a *App) {
	a.BrandChanges.Close()
	a.InfluencerChanges.Close()
	a.CustomerChanges.Close()
}

// InitProducer initializes kafka message producers
func InitProducer(a *App) {
	a.BrandFullProducer = kafka.NewSegmentioProducer(&kafka.SegmentioProducerOpts{
		Logger: a.Logger,
		Config: &a.Config.BrandFullProduceConfig,
	})

	a.InfluencerFullProducer = kafka.NewSegmentioProducer(&kafka.SegmentioProducerOpts{
		Logger: a.Logger,
		Config: &a.Config.InfluencerFullProducerConfig,
	})

}

// CloseProducer terminates all producer connections
func CloseProducer(a *App) {
	a.BrandFullProducer.Close()
	a.InfluencerFullProducer.Close()
}

func InitProcessor(a *App) {
	a.BrandProcessor = InitBrandProcessor(&BrandProcessorOpts{App: a, Logger: a.Logger})
	a.InfluencerProcessor = InitInfluencerProcessor(&InfluencerProcessorOpts{App: a, Logger: a.Logger})
	a.UserProcessor = InitUserProcessorOpts(&UserProcessorOpts{App: a, Logger: a.Logger})
}
