package app

import (
	"context"
	"go-app/server/kafka"
)

// InitConsumer initializes all kafka consumers
func InitConsumer(a *App) {
	ctx := context.TODO()

	a.UserChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.UserChangeConfig,
	})
	go a.UserChanges.Consume(ctx, a.UserProcessor.ProcessUserUpdate)

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
