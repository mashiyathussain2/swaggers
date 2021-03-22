package kafka

import (
	"context"
	"crypto/tls"
	"fmt"
	"go-app/server/config"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

// SegmentioKafkaImpl has kafka cluster config and connection instance
type SegmentioKafkaImpl struct {
	Config *config.KafkaConfig
	Conn   *kafka.Conn
}

// Close closes the connection
func (k *SegmentioKafkaImpl) Close() {
	k.Conn.Close()
}

// NewSegmentioKafka returns new segmentio kafka client instance
func NewSegmentioKafka(c *config.KafkaConfig) *SegmentioKafkaImpl {
	conn, err := kafka.Dial(c.BrokerDial, fmt.Sprintf("%s:%s", c.BrokerURL, c.BrokerPort))
	if err != nil {
		log.Fatalf("failed to establish kafka connection: %s", err)
		os.Exit(1)
	}
	defer conn.Close()
	controller, err := conn.Controller()
	if err != nil {
		log.Fatalf("failed while establishing connection to controller kafka: %s", err)
		os.Exit(1)
	}
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Fatalf("failed to establish connection to controller kafka: %s", err)
		os.Exit(1)
	}
	return &SegmentioKafkaImpl{Config: c, Conn: controllerConn}
}

// NewSegmentioKafkaDialer returns new segmentio kafka dialer instance
func NewSegmentioKafkaDialer(c *config.KafkaConfig) *kafka.Dialer {
	mechanism := plain.Mechanism{
		Username: c.Username,
		Password: c.Password,
	}
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		TLS:           &tls.Config{},
		SASLMechanism: mechanism,
	}
	return dialer
}

// SegmentioConsumer implements `Consumer` methods
type SegmentioConsumer struct {
	Reader *kafka.Reader
	Logger *zerolog.Logger
}

// SegmentioConsumerOpts contains args required to create SegmentioConsumer instance
type SegmentioConsumerOpts struct {
	Logger *zerolog.Logger
	Config *config.ListenerConfig
}

// NewSegmentioKafkaConsumer returns an instance of kafka segmentio consumer
func NewSegmentioKafkaConsumer(opts *SegmentioConsumerOpts) *SegmentioConsumer {
	s := SegmentioConsumer{Logger: opts.Logger}
	s.Init(opts.Config)
	return &s
}

// Init initialize kafka consumer group
func (cl *SegmentioConsumer) Init(c *config.ListenerConfig) {
	cl.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.Brokers,
		GroupID:  c.GroupID,
		Topic:    c.Topic,
		MaxBytes: 10e6, // 10MB
	})
}

// Consume consumes messages from kafka topic but does not commit them
func (cl *SegmentioConsumer) Consume(ctx context.Context, f func(Message)) {
	for {
		m, err := cl.Reader.FetchMessage(ctx)
		if err != nil {
			cl.Logger.Err(err).Msg("failed to fetch messages")
			break
		}
		f(m)
	}
}

// Commit commits an existing message
func (cl *SegmentioConsumer) Commit(ctx context.Context, m Message) {
	if err := cl.Reader.CommitMessages(ctx, m.(kafka.Message)); err != nil {
		cl.Logger.Err(err).Msg("failed to commit messages")
	}
}

// ConsumeAndCommit consumes a message and commits instantly
func (cl *SegmentioConsumer) ConsumeAndCommit(ctx context.Context, f func(Message)) {
	for {
		m, err := cl.Reader.ReadMessage(ctx)
		if err != nil {
			cl.Logger.Err(err).Msg("failed to fetch messages")
			break
		}
		f(m)
	}
}

// Close closes the consumer connection
func (cl *SegmentioConsumer) Close() {
	if err := cl.Reader.Close(); err != nil {
		cl.Logger.Err(err).Msg("failed to close reader")
	}
}
