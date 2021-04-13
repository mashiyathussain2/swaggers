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

	a.CatalogChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.CatalogChangeConfig,
	})
	go a.CatalogChanges.ConsumeAndCommit(ctx, a.CatalogProcessor.ProcessCatalogUpdate)

	a.CollectionCatalogChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.CollectionCatalogChangeConfig,
	})
	go a.CollectionCatalogChanges.ConsumeAndCommit(ctx, a.CollectionProcessor.ProcessCatalogUpdate)

	a.InventoryChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.InventoryChangeConfig,
	})
	go a.InventoryChanges.ConsumeAndCommit(ctx, a.CatalogProcessor.ProcessInventoryUpdate)

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

	go RunEvery(10*time.Second, a.Discount.CheckAndUpdateStatus)

}

// CloseConsumer close all consumer connections
func CloseConsumer(a *App) {
	a.ContentChanges.Close()
	a.CatalogChanges.Close()
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
}
