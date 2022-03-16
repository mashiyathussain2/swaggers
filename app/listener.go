package app

import (
	"context"
	"go-app/server/kafka"
	"time"
)

func RunEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}

// InitConsumer initializes consumers
func InitConsumer(a *App) {
	ctx := context.TODO()

	// a.CatalogChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
	// 	Logger: a.Logger,
	// 	Config: &a.Config.CatalogChangeConfig,
	// })

	// Catalog change consumer
	for i := 0; i < a.Config.CatalogChangeConfig.ConsumerCount; i++ {
		a.CatalogChanges = append(a.CatalogChanges, kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
			Logger: a.Logger,
			Config: &a.Config.CatalogChangeConfig,
		}))
	}
	for i := 0; i < a.Config.CatalogChangeConfig.ConsumerCount; i++ {
		for _, consumer := range a.CatalogChanges {
			go consumer.ConsumeAndCommit(ctx, a.CatalogProcessor.ProcessCatalogUpdate)
		}
	}

	a.InventoryChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.InventoryChangeConfig,
	})
	go a.InventoryChanges.ConsumeAndCommit(ctx, a.CatalogProcessor.ProcessInventoryUpdate)

	a.BrandChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.BrandChangeConfig,
	})
	go a.BrandChanges.ConsumeAndCommit(ctx, a.CatalogProcessor.ProcessBrandUpdate)

	a.DiscountChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.DiscountChangeConfig,
	})
	go a.DiscountChanges.ConsumeAndCommit(ctx, a.CatalogProcessor.ProcessDiscountUpdate)

	a.ContentChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.ContentChangeConfig,
	})
	go a.ContentChanges.ConsumeAndCommit(ctx, a.CatalogProcessor.ProcessCatalogContentUpdate)

	a.GroupChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.GroupChangeConfig,
	})
	go a.GroupChanges.ConsumeAndCommit(ctx, a.CatalogProcessor.ProcessGroupUpdate)

	a.CollectionChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.CollectionChangeConfig,
	})
	go a.CollectionChanges.ConsumeAndCommit(ctx, a.CollectionProcessor.ProcessCollectionUpdate)

	a.CollectionCatalogChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.CollectionCatalogChangeConfig,
	})
	go a.CollectionCatalogChanges.ConsumeAndCommit(ctx, a.CollectionProcessor.ProcessCatalogUpdate)

	a.ReviewChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.ReviewChangeConfig,
	})
	go a.ReviewChanges.ConsumeAndCommit(ctx, a.ReviewProcessor.ProcessReviewUpdate)

	go RunEvery(10*time.Second, a.Discount.CheckAndUpdateStatus)

}

// CloseConsumer close all consumer connections
func CloseConsumer(a *App) {
	a.ContentChanges.Close()
	// a.CatalogChanges.Close()
	for i := 0; i < a.Config.CatalogChangeConfig.ConsumerCount; i++ {
		a.CatalogChanges[i].Close()
	}
	a.BrandChanges.Close()
	a.CollectionCatalogChanges.Close()
	a.InventoryChanges.Close()
	a.DiscountChanges.Close()
	a.CollectionChanges.Close()
	a.GroupChanges.Close()
}

// InitProducer initializes kafka message producers
func InitProducer(a *App) {
	a.CatalogFullProducer = kafka.NewSegmentioProducer(&kafka.SegmentioProducerOpts{
		Logger: a.Logger,
		Config: &a.Config.CatalogFullProducerConfig,
	})

	a.CollectionFullProducer = kafka.NewSegmentioProducer(&kafka.SegmentioProducerOpts{
		Logger: a.Logger,
		Config: &a.Config.CollectionFullProducerConfig,
	})

	a.ReviewFullProducer = kafka.NewSegmentioProducer(&kafka.SegmentioProducerOpts{
		Logger: a.Logger,
		Config: &a.Config.ReviewFullProducerConfig,
	})
}

// CloseProducer terminates all producer connections
func CloseProducer(a *App) {
	a.CatalogFullProducer.Close()
	a.CollectionFullProducer.Close()
}

func InitProcessor(a *App) {
	a.CatalogProcessor = InitCatalogProcessor(&CatalogProcessorOpts{
		App:    a,
		Logger: a.Logger,
	})

	a.CollectionProcessor = InitCollectionProcessor(&CollectionProcessorOpts{
		App:    a,
		Logger: a.Logger,
	})

	a.ReviewProcessor = InitReviewProcessor(&ReviewProcessorOpts{
		App:    a,
		Logger: a.Logger,
	})
}
