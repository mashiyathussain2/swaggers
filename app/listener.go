package app

import (
	"context"
	"go-app/server/kafka"
)

// InitConsumer initializes consumers
func InitConsumer(a *App) {
	ctx := context.Background()

	a.LiveComments = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.LiveCommentChangesConfig,
	})
	go a.LiveComments.Consume(ctx, a.Live.ConsumeComment)

	a.BrandChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.BrandChangesConfig,
	})
	go a.BrandChanges.ConsumeAndCommit(ctx, a.ContentUpdateProcessor.ProcessBrandMessage)

	a.InfluencerChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.InfluencerChangesConfig,
	})
	go a.InfluencerChanges.ConsumeAndCommit(ctx, a.ContentUpdateProcessor.ProcessInfluencerMessage)

	a.CatalogChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.CatalogChangesConfig,
	})
	go a.CatalogChanges.ConsumeAndCommit(ctx, a.ContentUpdateProcessor.ProcessCatalogMessage)

	a.ContentChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.ContentChangesConfig,
	})
	go a.ContentChanges.ConsumeAndCommit(ctx, a.ContentUpdateProcessor.ProcessContentMessage)

}

// CloseConsumer close all consumer connections
func CloseConsumer(a *App) {
	a.LiveComments.Close()
	a.BrandChanges.Close()
	a.InfluencerChanges.Close()
	a.CatalogChanges.Close()
}

func InitProducer(a *App) {
	a.LiveCommentProducer = kafka.NewSegmentioProducer(&kafka.SegmentioProducerOpts{
		Logger: a.Logger,
		Config: &a.Config.LiveCommentProducerConfig,
	})

	a.ContentFullProducer = kafka.NewSegmentioProducer(&kafka.SegmentioProducerOpts{
		Logger: a.Logger,
		Config: &a.Config.ContentFullProducerConfig,
	})

}

func CloseProducer(a *App) {
	a.LiveCommentProducer.Close()
	a.ContentFullProducer.Close()
}

func InitProcessor(a *App) {
	a.ContentUpdateProcessor = InitContentUpdateProcessor(&ContentUpdateProcessorOpts{
		App:    a,
		Logger: a.Logger,
	})
}
