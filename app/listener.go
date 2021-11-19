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

	a.LikeChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.LikeChangeConfig,
	})
	go a.LikeChanges.ConsumeAndCommit(ctx, a.ContentUpdateProcessor.ProcessLike)

	a.ViewChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.ViewChangeConfig,
	})
	go a.ViewChanges.ConsumeAndCommit(ctx, a.ContentUpdateProcessor.ProcessView)

	a.CommentChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.CommentChangeConfig,
	})
	go a.CommentChanges.ConsumeAndCommit(ctx, a.ContentUpdateProcessor.ProcessComment)

	a.PebbleSeriesConsumer = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.SeriesConsumerConfig,
	})
	go a.PebbleSeriesConsumer.ConsumeAndCommit(ctx, a.ContentUpdateProcessor.ProcessSeriesMessage)

	a.PebbleCollectionConsumer = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.CollectionConsumerConfig,
	})
	go a.PebbleCollectionConsumer.ConsumeAndCommit(ctx, a.ContentUpdateProcessor.ProcessCollectionMessage)

	a.PebbleStatusChangeForSeries = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.PebbleStatusChangeForSeriesConfig,
	})
	go a.PebbleStatusChangeForSeries.ConsumeAndCommit(ctx, a.ContentUpdateProcessor.ProcessContentMessageForSeries)

	a.LikeChangeForSeries = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.LikeChangeForSeriesConfig,
	})
	go a.LikeChangeForSeries.ConsumeAndCommit(ctx, a.ContentUpdateProcessor.ProcessLikeForSeries)
}

// CloseConsumer close all consumer connections
func CloseConsumer(a *App) {
	a.LiveComments.Close()
	a.BrandChanges.Close()
	a.InfluencerChanges.Close()
	a.LikeChanges.Close()
	a.CommentChanges.Close()
	a.ViewChanges.Close()
	a.CatalogChanges.Close()
	a.ContentChanges.Close()
	a.PebbleSeriesConsumer.Close()
	a.PebbleCollectionConsumer.Close()
	a.LikeChangeForSeries.Close()
	a.PebbleStatusChangeForSeries.Close()
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
	a.PebbleSeriesProducer = kafka.NewSegmentioProducer(&kafka.SegmentioProducerOpts{
		Logger: a.Logger,
		Config: &a.Config.SeriesFullProducerConfig,
	})
	a.PebbleCollectionProducer = kafka.NewSegmentioProducer(&kafka.SegmentioProducerOpts{
		Logger: a.Logger,
		Config: &a.Config.CollectionFullProducerConfig,
	})
}

func CloseProducer(a *App) {
	a.LiveCommentProducer.Close()
	a.ContentFullProducer.Close()
	a.PebbleSeriesProducer.Close()
	a.PebbleCollectionProducer.Close()
}

func InitProcessor(a *App) {
	a.ContentUpdateProcessor = InitContentUpdateProcessor(&ContentUpdateProcessorOpts{
		App:    a,
		Logger: a.Logger,
	})
}
