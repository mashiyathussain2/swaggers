package app

import (
	"context"
	"go-app/server/kafka"
)

// InitConsumer initializes consumers
func InitConsumer(a *App) {
	ctx := context.Background()

	// Captures comments from comments collection and specifically publishes live comments to AWS IVS Channel for broadcast
	a.LiveComments = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.LiveCommentChangesConfig,
	})
	go a.LiveComments.Consume(ctx, a.Live.ConsumeComment)

	// Captures brand related changes such as name, logo etc and syncs them with the content collections.
	a.BrandChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.BrandChangesConfig,
	})
	go a.BrandChanges.ConsumeAndCommit(ctx, a.ContentUpdateProcessor.ProcessBrandMessage)

	// Captures influencer related changes such as name, logo etc and syncs them with the content collections.
	a.InfluencerChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.InfluencerChangesConfig,
	})
	go a.InfluencerChanges.ConsumeAndCommit(ctx, a.ContentUpdateProcessor.ProcessInfluencerMessage)

	// Captures catalog related changes such as name, description, price etc and syncs them with the content collections.
	a.CatalogChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.CatalogChangesConfig,
	})
	go a.CatalogChanges.ConsumeAndCommit(ctx, a.ContentUpdateProcessor.ProcessCatalogMessage)

	// Content change consumer
	for i := 0; i < a.Config.ContentChangesConfig.ConsumerCount; i++ {
		a.ContentChanges = append(a.ContentChanges, kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
			Logger: a.Logger,
			Config: &a.Config.ContentChangesConfig,
		}))
	}
	for i := 0; i < a.Config.ContentChangesConfig.ConsumerCount; i++ {
		for _, consumer := range a.ContentChanges {
			go consumer.ConsumeAndCommit(ctx, a.ContentUpdateProcessor.ProcessContentMessage)
		}
	}

	// Like consumers
	for i := 0; i < a.Config.LikeChangeConfig.ConsumerCount; i++ {
		a.LikeChanges = append(a.LikeChanges, kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
			Logger: a.Logger,
			Config: &a.Config.LikeChangeConfig,
		}))
	}
	for i := 0; i < a.Config.LikeChangeConfig.ConsumerCount; i++ {
		for _, consumer := range a.LikeChanges {
			go consumer.ConsumeAndCommit(ctx, a.ContentUpdateProcessor.ProcessLike)
		}
	}

	// View consumers
	for i := 0; i < a.Config.ViewChangeConfig.ConsumerCount; i++ {
		a.ViewChanges = append(a.ViewChanges, kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
			Logger: a.Logger,
			Config: &a.Config.ViewChangeConfig,
		}))
	}
	for i := 0; i < a.Config.ViewChangeConfig.ConsumerCount; i++ {
		for _, consumer := range a.ViewChanges {
			go consumer.ConsumeAndCommit(ctx, a.ContentUpdateProcessor.ProcessView)
		}
	}

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

}

// CloseConsumer close all consumer connections
func CloseConsumer(a *App) {
	a.LiveComments.Close()
	a.BrandChanges.Close()
	a.InfluencerChanges.Close()
	a.CommentChanges.Close()

	for i := 0; i < a.Config.ViewChangeConfig.ConsumerCount; i++ {
		a.ViewChanges[i].Close()
	}
	for i := 0; i < a.Config.LikeChangeConfig.ConsumerCount; i++ {
		a.LikeChanges[i].Close()
	}

	for i := 0; i < a.Config.ContentChangesConfig.ConsumerCount; i++ {
		a.ContentChanges[i].Close()
	}

	a.CatalogChanges.Close()
	a.PebbleSeriesConsumer.Close()
	a.PebbleCollectionConsumer.Close()
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
	a.LikeProducer = kafka.NewSegmentioProducer(&kafka.SegmentioProducerOpts{
		Logger: a.Logger,
		Config: &a.Config.LikeProducerConfig,
	})

	a.ViewProducer = kafka.NewSegmentioProducer(&kafka.SegmentioProducerOpts{
		Logger: a.Logger,
		Config: &a.Config.ViewProducerConfig,
	})
}

func CloseProducer(a *App) {
	a.LiveCommentProducer.Close()
	a.ContentFullProducer.Close()
	a.PebbleSeriesProducer.Close()
	a.PebbleCollectionProducer.Close()
	a.LikeProducer.Close()
	a.ViewProducer.Close()
}

func InitProcessor(a *App) {
	a.ContentUpdateProcessor = InitContentUpdateProcessor(&ContentUpdateProcessorOpts{
		App:    a,
		Logger: a.Logger,
	})
}
