package pulsarutil

import (
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

// ClientBuilder provides a builder pattern for creating Pulsar clients
type ClientBuilder struct {
	config *Config
}

// NewClientBuilder creates a new client builder
func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		config: DefaultConfig(),
	}
}

// ServiceURL sets the Pulsar service URL
func (cb *ClientBuilder) ServiceURL(url string) *ClientBuilder {
	cb.config.ServiceURL = url
	return cb
}

// ConnectionTimeout sets the connection timeout
func (cb *ClientBuilder) ConnectionTimeout(timeout time.Duration) *ClientBuilder {
	cb.config.ConnectionTimeout = timeout
	return cb
}

// OperationTimeout sets the operation timeout
func (cb *ClientBuilder) OperationTimeout(timeout time.Duration) *ClientBuilder {
	cb.config.OperationTimeout = timeout
	return cb
}

// MaxConnectionsPerBroker sets the maximum connections per broker
func (cb *ClientBuilder) MaxConnectionsPerBroker(max int) *ClientBuilder {
	cb.config.MaxConnectionsPerBroker = max
	return cb
}

// Authentication sets the authentication method
func (cb *ClientBuilder) Authentication(auth pulsar.Authentication) *ClientBuilder {
	cb.config.Authentication = auth
	return cb
}

// TLSAllowInsecureConnection allows insecure TLS connections
func (cb *ClientBuilder) TLSAllowInsecureConnection(allow bool) *ClientBuilder {
	cb.config.TLSAllowInsecureConnection = allow
	return cb
}

// TLSTrustCertsFilePath sets the TLS trust certificates file path
func (cb *ClientBuilder) TLSTrustCertsFilePath(path string) *ClientBuilder {
	cb.config.TLSTrustCertsFilePath = path
	return cb
}

// Build creates the Pulsar client
func (cb *ClientBuilder) Build() (*Client, error) {
	return NewClient(cb.config)
}

// ProducerBuilder provides a builder pattern for creating producers
type ProducerBuilder struct {
	client *Client
	config *ProducerConfig
}

// NewProducerBuilder creates a new producer builder
func (c *Client) NewProducerBuilder(topic string) *ProducerBuilder {
	return &ProducerBuilder{
		client: c,
		config: DefaultProducerConfig(topic),
	}
}

// ProducerName sets the producer name
func (pb *ProducerBuilder) ProducerName(name string) *ProducerBuilder {
	pb.config.ProducerName = name
	return pb
}

// SendTimeout sets the send timeout
func (pb *ProducerBuilder) SendTimeout(timeout time.Duration) *ProducerBuilder {
	pb.config.SendTimeout = timeout
	return pb
}

// DisableBatching disables message batching
func (pb *ProducerBuilder) DisableBatching(disable bool) *ProducerBuilder {
	pb.config.DisableBatching = disable
	return pb
}

// BatchingMaxPublishDelay sets the maximum batching publish delay
func (pb *ProducerBuilder) BatchingMaxPublishDelay(delay time.Duration) *ProducerBuilder {
	pb.config.BatchingMaxPublishDelay = delay
	return pb
}

// BatchingMaxMessages sets the maximum number of messages in a batch
func (pb *ProducerBuilder) BatchingMaxMessages(max uint) *ProducerBuilder {
	pb.config.BatchingMaxMessages = max
	return pb
}

// BatchingMaxSize sets the maximum batch size
func (pb *ProducerBuilder) BatchingMaxSize(size uint) *ProducerBuilder {
	pb.config.BatchingMaxSize = size
	return pb
}

// MaxPendingMessages sets the maximum pending messages
func (pb *ProducerBuilder) MaxPendingMessages(max int) *ProducerBuilder {
	pb.config.MaxPendingMessages = max
	return pb
}

// HashingScheme sets the hashing scheme
func (pb *ProducerBuilder) HashingScheme(scheme pulsar.HashingScheme) *ProducerBuilder {
	pb.config.HashingScheme = scheme
	return pb
}

// CompressionType sets the compression type
func (pb *ProducerBuilder) CompressionType(compression pulsar.CompressionType) *ProducerBuilder {
	pb.config.CompressionType = compression
	return pb
}

// CompressionLevel sets the compression level
func (pb *ProducerBuilder) CompressionLevel(level pulsar.CompressionLevel) *ProducerBuilder {
	pb.config.CompressionLevel = level
	return pb
}

// DisableReplication disables message replication
func (pb *ProducerBuilder) DisableReplication(disable bool) *ProducerBuilder {
	pb.config.DisableReplication = disable
	return pb
}

// Build creates the concurrent producer
func (pb *ProducerBuilder) Build(workerCount int) (*ConcurrentProducer, error) {
	return pb.client.NewConcurrentProducer(pb.config, workerCount)
}

// ConsumerBuilder provides a builder pattern for creating consumers
type ConsumerBuilder struct {
	client *Client
	config *ConsumerConfig
}

// NewConsumerBuilder creates a new consumer builder
func (c *Client) NewConsumerBuilder(topics []string, subscription string) *ConsumerBuilder {
	return &ConsumerBuilder{
		client: c,
		config: DefaultConsumerConfig(topics, subscription),
	}
}

// SubscriptionType sets the subscription type
func (cb *ConsumerBuilder) SubscriptionType(subType pulsar.SubscriptionType) *ConsumerBuilder {
	cb.config.SubscriptionType = subType
	return cb
}

// SubscriptionInitialPosition sets the initial position for the subscription
func (cb *ConsumerBuilder) SubscriptionInitialPosition(position pulsar.SubscriptionInitialPosition) *ConsumerBuilder {
	cb.config.SubscriptionInitialPosition = position
	return cb
}

// MessageChannel sets the message channel buffer size
func (cb *ConsumerBuilder) MessageChannel(size int) *ConsumerBuilder {
	cb.config.MessageChannel = size
	return cb
}

// ReceiverQueueSize sets the receiver queue size
func (cb *ConsumerBuilder) ReceiverQueueSize(size int) *ConsumerBuilder {
	cb.config.ReceiverQueueSize = size
	return cb
}

// MaxUnackedMessages sets the maximum unacked messages
func (cb *ConsumerBuilder) MaxUnackedMessages(max int) *ConsumerBuilder {
	cb.config.MaxUnackedMessages = max
	return cb
}

// AckTimeout sets the ack timeout
func (cb *ConsumerBuilder) AckTimeout(timeout time.Duration) *ConsumerBuilder {
	cb.config.AckTimeout = timeout
	return cb
}

// NackRedeliveryDelay sets the nack redelivery delay
func (cb *ConsumerBuilder) NackRedeliveryDelay(delay time.Duration) *ConsumerBuilder {
	cb.config.NackRedeliveryDelay = delay
	return cb
}

// RetryEnable enables retry functionality
func (cb *ConsumerBuilder) RetryEnable(enable bool) *ConsumerBuilder {
	cb.config.RetryEnable = enable
	return cb
}

// DLQ sets the Dead Letter Queue policy
func (cb *ConsumerBuilder) DLQ(dlq *pulsar.DLQPolicy) *ConsumerBuilder {
	cb.config.DLQ = dlq
	return cb
}

// ReadCompacted enables reading compacted messages
func (cb *ConsumerBuilder) ReadCompacted(enable bool) *ConsumerBuilder {
	cb.config.ReadCompacted = enable
	return cb
}

// ReplicateSubscriptionState enables subscription state replication
func (cb *ConsumerBuilder) ReplicateSubscriptionState(enable bool) *ConsumerBuilder {
	cb.config.ReplicateSubscriptionState = enable
	return cb
}

// Properties sets consumer properties
func (cb *ConsumerBuilder) Properties(props map[string]string) *ConsumerBuilder {
	cb.config.Properties = props
	return cb
}

// Build creates the concurrent consumer
func (cb *ConsumerBuilder) Build(handler MessageHandler, autoAck bool, workerCount int) (*ConcurrentConsumer, error) {
	return cb.client.NewConcurrentConsumer(cb.config, handler, autoAck, workerCount)
}

// BuildChannelBased creates a concurrent consumer that uses channels instead of handlers
func (cb *ConsumerBuilder) BuildChannelBased(workerCount int) (*ConcurrentConsumer, error) {
	return cb.client.NewConcurrentConsumer(cb.config, nil, false, workerCount)
}
