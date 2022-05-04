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
	go a.DiscountChanges.ConsumeAndCommit(ctx, a.CartProcessor.ProcessDiscountUpdate)

	a.CatalogChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.CatalogChangeConfig,
	})
	go a.CatalogChanges.ConsumeAndCommit(ctx, a.CartProcessor.ProcessCatalogUpdate)

	a.InventoryChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.InventoryChangeConfig,
	})
	go a.InventoryChanges.ConsumeAndCommit(ctx, a.CartProcessor.ProcessInventoryUpdate)

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

	a.CouponChanges = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.CouponChangeConfig,
	})
	go a.CouponChanges.ConsumeAndCommit(ctx, a.CartProcessor.ProcessCouponUpdate)

	a.CommissionOrderListener = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.CommissionOrderListenerConfig,
	})
	go a.CommissionOrderListener.ConsumeAndCommit(ctx, a.InfluencerProcessor.InfluencerCommissionUpdate)

	a.GenerateCommissionInvoiceListner = kafka.NewSegmentioKafkaConsumer(&kafka.SegmentioConsumerOpts{
		Logger: a.Logger,
		Config: &a.Config.GenerateCommissionInvoiceListnerConfig,
	})
	go a.CommissionOrderListener.ConsumeAndCommit(ctx, a.InfluencerProcessor.GenerateCommissionInvoice)

}

// CloseConsumer terminates all consumer connections
func CloseConsumer(a *App) {
	a.BrandChanges.Close()
	a.InfluencerChanges.Close()
	a.CustomerChanges.Close()
	a.DiscountChanges.Close()
	// a.InfluencerChanges.Close()
	a.CatalogChanges.Close()
	a.CouponChanges.Close()
	a.GenerateCommissionInvoiceListner.Close()
	a.CommissionOrderListener.Close()
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
	a.GenerateCommissionInvoiceProducer = kafka.NewSegmentioProducer(&kafka.SegmentioProducerOpts{
		Logger: a.Logger,
		Config: &a.Config.GenerateCommissionInvoiceProducerConfig,
	})

}

// CloseProducer terminates all producer connections
func CloseProducer(a *App) {
	a.BrandFullProducer.Close()
	a.InfluencerFullProducer.Close()
	a.GenerateCommissionInvoiceProducer.Close()
}

func InitProcessor(a *App) {
	a.BrandProcessor = InitBrandProcessor(&BrandProcessorOpts{App: a, Logger: a.Logger})
	a.InfluencerProcessor = InitInfluencerProcessor(&InfluencerProcessorOpts{App: a, Logger: a.Logger})
	a.UserProcessor = InitUserProcessorOpts(&UserProcessorOpts{App: a, Logger: a.Logger})
	a.CartProcessor = InitCartProcessorOpts(&CartProcessorOpts{App: a, Logger: a.Logger})
}
