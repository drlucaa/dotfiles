package pulsarutil

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

// ConsumerMessage represents a received message
type ConsumerMessage struct {
	Topic           string
	Subscription    string
	Payload         []byte
	Key             string
	Properties      map[string]string
	PublishTime     time.Time
	EventTime       time.Time
	MessageID       pulsar.MessageID
	RedeliveryCount uint32
	Message         pulsar.Message // Reference to original message for acking
}

// MessageHandler defines the interface for handling received messages
type MessageHandler func(ctx context.Context, msg *ConsumerMessage) error

// ConsumerConfig holds configuration for consumers
type ConsumerConfig struct {
	Topics                      []string
	Subscription                string
	SubscriptionType            pulsar.SubscriptionType
	SubscriptionInitialPosition pulsar.SubscriptionInitialPosition
	MessageChannel              int
	ReceiverQueueSize           int
	MaxUnackedMessages          int
	AckTimeout                  time.Duration
	NackRedeliveryDelay         time.Duration
	RetryEnable                 bool
	DLQ                         *pulsar.DLQPolicy
	ReadCompacted               bool
	ReplicateSubscriptionState  bool
	Properties                  map[string]string
}

// ConcurrentConsumer provides a channel-based interface for concurrent message consumption
type ConcurrentConsumer struct {
	client      *Client
	config      *ConsumerConfig
	consumers   map[string]pulsar.Consumer
	mu          sync.RWMutex
	messageCh   chan *ConsumerMessage
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	handler     MessageHandler
	autoAck     bool
	workerCount int
}

// DefaultConsumerConfig returns a default consumer configuration
func DefaultConsumerConfig(topics []string, subscription string) *ConsumerConfig {
	return &ConsumerConfig{
		Topics:                      topics,
		Subscription:                subscription,
		SubscriptionType:            pulsar.Shared,
		SubscriptionInitialPosition: pulsar.SubscriptionPositionLatest,
		MessageChannel:              1000,
		ReceiverQueueSize:           1000,
		MaxUnackedMessages:          1000,
		AckTimeout:                  30 * time.Second,
		NackRedeliveryDelay:         1 * time.Minute,
		RetryEnable:                 true,
		ReadCompacted:               false,
		ReplicateSubscriptionState:  false,
		Properties:                  make(map[string]string),
	}
}

// NewConcurrentConsumer creates a new concurrent consumer
func (c *Client) NewConcurrentConsumer(config *ConsumerConfig, handler MessageHandler, autoAck bool, workerCount int) (*ConcurrentConsumer, error) {
	if config == nil {
		return nil, fmt.Errorf("consumer config cannot be nil")
	}

	if len(config.Topics) == 0 {
		return nil, fmt.Errorf("at least one topic must be specified")
	}

	if config.Subscription == "" {
		return nil, fmt.Errorf("subscription name cannot be empty")
	}

	if workerCount <= 0 {
		workerCount = 1
	}

	ctx, cancel := context.WithCancel(context.Background())

	consumer := &ConcurrentConsumer{
		client:      c,
		config:      config,
		consumers:   make(map[string]pulsar.Consumer),
		messageCh:   make(chan *ConsumerMessage, config.MessageChannel),
		ctx:         ctx,
		cancel:      cancel,
		handler:     handler,
		autoAck:     autoAck,
		workerCount: workerCount,
	}

	// Create consumers for each topic
	for _, topic := range config.Topics {
		pulsarConsumer, err := consumer.createConsumer(topic)
		if err != nil {
			consumer.Close()
			return nil, fmt.Errorf("failed to create consumer for topic %s: %w", topic, err)
		}
		consumer.consumers[topic] = pulsarConsumer
	}

	// Start message receiving goroutines
	for topic, pulsarConsumer := range consumer.consumers {
		consumer.wg.Add(1)
		go consumer.receiveMessages(topic, pulsarConsumer)
	}

	// Start worker goroutines if handler is provided
	if handler != nil {
		for i := 0; i < workerCount; i++ {
			consumer.wg.Add(1)
			go consumer.worker(i)
		}
	}

	c.logger.Infof("Created concurrent consumer with %d workers for topics: %v, subscription: %s",
		workerCount, config.Topics, config.Subscription)
	return consumer, nil
}

// createConsumer creates a Pulsar consumer for a specific topic
func (cc *ConcurrentConsumer) createConsumer(topic string) (pulsar.Consumer, error) {
	consumerOptions := pulsar.ConsumerOptions{
		Topic:                       topic,
		SubscriptionName:            cc.config.Subscription,
		Type:                        cc.config.SubscriptionType,
		SubscriptionInitialPosition: cc.config.SubscriptionInitialPosition,
		ReceiverQueueSize:           cc.config.ReceiverQueueSize,
		NackRedeliveryDelay:         cc.config.NackRedeliveryDelay,
		RetryEnable:                 cc.config.RetryEnable,
		ReadCompacted:               cc.config.ReadCompacted,
		ReplicateSubscriptionState:  cc.config.ReplicateSubscriptionState,
		Properties:                  cc.config.Properties,
	}

	if cc.config.DLQ != nil {
		consumerOptions.DLQ = cc.config.DLQ
	}

	consumer, err := cc.client.client.Subscribe(consumerOptions)
	if err != nil {
		return nil, err
	}

	cc.client.logger.Infof("Created consumer for topic: %s, subscription: %s", topic, cc.config.Subscription)
	return consumer, nil
}

// receiveMessages continuously receives messages from a Pulsar consumer
func (cc *ConcurrentConsumer) receiveMessages(topic string, consumer pulsar.Consumer) {
	defer cc.wg.Done()
	cc.client.logger.Debugf("Message receiver started for topic: %s", topic)

	for {
		select {
		case <-cc.ctx.Done():
			cc.client.logger.Debugf("Message receiver stopping for topic: %s", topic)
			return
		default:
			msg, err := consumer.Receive(cc.ctx)
			if err != nil {
				if cc.ctx.Err() != nil {
					return // Context cancelled
				}
				cc.client.logger.Errorf("Error receiving message from topic %s: %v", topic, err)
				continue
			}

			consumerMsg := &ConsumerMessage{
				Topic:           topic,
				Subscription:    cc.config.Subscription,
				Payload:         msg.Payload(),
				Key:             msg.Key(),
				Properties:      msg.Properties(),
				PublishTime:     msg.PublishTime(),
				EventTime:       msg.EventTime(),
				MessageID:       msg.ID(),
				RedeliveryCount: msg.RedeliveryCount(),
				Message:         msg,
			}

			// Send message to channel, but don't block indefinitely
			select {
			case cc.messageCh <- consumerMsg:
				cc.client.logger.Debugf("Message received from topic %s: %s", topic, msg.ID())
			case <-cc.ctx.Done():
				// Nack the message since we're shutting down
				consumer.Nack(msg)
				return
			default:
				cc.client.logger.Warn("Message channel full, nacking message")
				consumer.Nack(msg)
			}
		}
	}
}

// worker processes messages from the message channel using the handler
func (cc *ConcurrentConsumer) worker(workerID int) {
	defer cc.wg.Done()
	cc.client.logger.Debugf("Consumer worker %d started", workerID)

	for {
		select {
		case <-cc.ctx.Done():
			cc.client.logger.Debugf("Consumer worker %d stopping", workerID)
			return
		case msg, ok := <-cc.messageCh:
			if !ok {
				cc.client.logger.Debugf("Consumer worker %d: message channel closed", workerID)
				return
			}

			cc.processMessage(msg)
		}
	}
}

// processMessage handles a single message using the configured handler
func (cc *ConcurrentConsumer) processMessage(msg *ConsumerMessage) {
	consumer := cc.consumers[msg.Topic]
	if consumer == nil {
		cc.client.logger.Errorf("No consumer found for topic: %s", msg.Topic)
		return
	}

	if cc.handler != nil {
		err := cc.handler(cc.ctx, msg)
		if err != nil {
			cc.client.logger.Errorf("Error processing message %s from topic %s: %v",
				msg.MessageID, msg.Topic, err)
			consumer.Nack(msg.Message)
		} else {
			if cc.autoAck {
				cc.client.logger.Debugf("Auto-acking message %s from topic %s",
					msg.MessageID, msg.Topic)
				consumer.Ack(msg.Message)
			}
		}
	}
}

// MessageChannel returns the channel for receiving messages (when no handler is set)
func (cc *ConcurrentConsumer) MessageChannel() <-chan *ConsumerMessage {
	return cc.messageCh
}

// Ack acknowledges a message
func (cc *ConcurrentConsumer) Ack(msg *ConsumerMessage) error {
	consumer := cc.consumers[msg.Topic]
	if consumer == nil {
		return fmt.Errorf("no consumer found for topic: %s", msg.Topic)
	}
	consumer.Ack(msg.Message)
	cc.client.logger.Debugf("Acknowledged message %s from topic %s", msg.MessageID, msg.Topic)
	return nil
}

// Nack negatively acknowledges a message
func (cc *ConcurrentConsumer) Nack(msg *ConsumerMessage) error {
	consumer := cc.consumers[msg.Topic]
	if consumer == nil {
		return fmt.Errorf("no consumer found for topic: %s", msg.Topic)
	}
	consumer.Nack(msg.Message)
	cc.client.logger.Debugf("Nacked message %s from topic %s", msg.MessageID, msg.Topic)
	return nil
}

// AckID acknowledges a message by its ID
func (cc *ConcurrentConsumer) AckID(topic string, msgID pulsar.MessageID) error {
	consumer := cc.consumers[topic]
	if consumer == nil {
		return fmt.Errorf("no consumer found for topic: %s", topic)
	}
	consumer.AckID(msgID)
	cc.client.logger.Debugf("Acknowledged message ID %s from topic %s", msgID, topic)
	return nil
}

// Close closes the concurrent consumer
func (cc *ConcurrentConsumer) Close() error {
	cc.cancel()

	// Wait for all goroutines to finish
	cc.wg.Wait()

	// Close all consumers
	cc.mu.Lock()
	defer cc.mu.Unlock()

	for topic, consumer := range cc.consumers {
		consumer.Close()
		cc.client.logger.Infof("Closed consumer for topic: %s", topic)
	}

	close(cc.messageCh)
	cc.client.logger.Info("Concurrent consumer closed")
	return nil
}

// Stats returns basic statistics about the consumer
func (cc *ConcurrentConsumer) Stats() map[string]interface{} {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	return map[string]interface{}{
		"topics":        len(cc.consumers),
		"workers":       cc.workerCount,
		"message_queue": len(cc.messageCh),
		"auto_ack":      cc.autoAck,
	}
}

// Seek seeks to a specific message ID or timestamp for all consumers
func (cc *ConcurrentConsumer) Seek(msgID pulsar.MessageID) error {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	for topic, consumer := range cc.consumers {
		if err := consumer.Seek(msgID); err != nil {
			return fmt.Errorf("failed to seek consumer for topic %s: %w", topic, err)
		}
		cc.client.logger.Infof("Seeked to message ID %s for topic %s", msgID, topic)
	}
	return nil
}

// SeekByTime seeks to a specific timestamp for all consumers
func (cc *ConcurrentConsumer) SeekByTime(timestamp time.Time) error {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	for topic, consumer := range cc.consumers {
		if err := consumer.SeekByTime(timestamp); err != nil {
			return fmt.Errorf("failed to seek consumer for topic %s to time %v: %w", topic, timestamp, err)
		}
		cc.client.logger.Infof("Seeked to timestamp %v for topic %s", timestamp, topic)
	}
	return nil
}
