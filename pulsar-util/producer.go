package pulsarutil

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

// ProducerMessage represents a message to be sent
type ProducerMessage struct {
	Topic               string
	Payload             []byte
	Key                 string
	Properties          map[string]string
	EventTime           time.Time
	ReplicationClusters []string
	DisableReplication  bool
}

// ProducerResult represents the result of a send operation
type ProducerResult struct {
	MessageID   pulsar.MessageID
	Message     *ProducerMessage
	Error       error
	PublishTime time.Time
}

// ProducerConfig holds configuration for producers
type ProducerConfig struct {
	Topic                   string
	ProducerName            string
	SendTimeout             time.Duration
	DisableBatching         bool
	BatchingMaxPublishDelay time.Duration
	BatchingMaxMessages     uint
	BatchingMaxSize         uint
	MaxPendingMessages      int
	HashingScheme           pulsar.HashingScheme
	CompressionType         pulsar.CompressionType
	CompressionLevel        pulsar.CompressionLevel
	DisableReplication      bool
	Partitions              string // auto, single, custom
}

// ConcurrentProducer provides a channel-based interface for concurrent message publishing
type ConcurrentProducer struct {
	client      *Client
	config      *ProducerConfig
	producers   map[string]pulsar.Producer
	mu          sync.RWMutex
	inputCh     chan *ProducerMessage
	resultCh    chan *ProducerResult
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	workerCount int
}

// DefaultProducerConfig returns a default producer configuration
func DefaultProducerConfig(topic string) *ProducerConfig {
	return &ProducerConfig{
		Topic:                   topic,
		SendTimeout:             30 * time.Second,
		DisableBatching:         false,
		BatchingMaxPublishDelay: 10 * time.Millisecond,
		BatchingMaxMessages:     1000,
		BatchingMaxSize:         128 * 1024, // 128KB
		MaxPendingMessages:      1000,
		HashingScheme:           pulsar.JavaStringHash,
		CompressionType:         pulsar.LZ4,
		CompressionLevel:        pulsar.Default,
		DisableReplication:      false,
		Partitions:              "auto",
	}
}

// NewConcurrentProducer creates a new concurrent producer
func (c *Client) NewConcurrentProducer(config *ProducerConfig, workerCount int) (*ConcurrentProducer, error) {
	if config == nil {
		return nil, fmt.Errorf("producer config cannot be nil")
	}

	if workerCount <= 0 {
		workerCount = 1
	}

	ctx, cancel := context.WithCancel(context.Background())

	producer := &ConcurrentProducer{
		client:      c,
		config:      config,
		producers:   make(map[string]pulsar.Producer),
		inputCh:     make(chan *ProducerMessage, config.MaxPendingMessages),
		resultCh:    make(chan *ProducerResult, config.MaxPendingMessages),
		ctx:         ctx,
		cancel:      cancel,
		workerCount: workerCount,
	}

	// Start worker goroutines
	for i := 0; i < workerCount; i++ {
		producer.wg.Add(1)
		go producer.worker(i)
	}

	c.logger.Infof("Created concurrent producer with %d workers for topic: %s", workerCount, config.Topic)
	return producer, nil
}

// worker processes messages from the input channel
func (cp *ConcurrentProducer) worker(workerID int) {
	defer cp.wg.Done()
	cp.client.logger.Debugf("Producer worker %d started", workerID)

	for {
		select {
		case <-cp.ctx.Done():
			cp.client.logger.Debugf("Producer worker %d stopping", workerID)
			return
		case msg, ok := <-cp.inputCh:
			if !ok {
				cp.client.logger.Debugf("Producer worker %d: input channel closed", workerID)
				return
			}

			result := cp.processMessage(msg)

			// Send result back, but don't block if result channel is full
			select {
			case cp.resultCh <- result:
			case <-cp.ctx.Done():
				return
			default:
				cp.client.logger.Warn("Result channel full, dropping result")
			}
		}
	}
}

// processMessage handles the actual message sending
func (cp *ConcurrentProducer) processMessage(msg *ProducerMessage) *ProducerResult {
	topic := msg.Topic
	if topic == "" {
		topic = cp.config.Topic
	}

	producer, err := cp.getOrCreateProducer(topic)
	if err != nil {
		return &ProducerResult{
			Message: msg,
			Error:   fmt.Errorf("failed to get producer for topic %s: %w", topic, err),
		}
	}

	pulsarMsg := &pulsar.ProducerMessage{
		Payload:    msg.Payload,
		Key:        msg.Key,
		Properties: msg.Properties,
	}

	if !msg.EventTime.IsZero() {
		pulsarMsg.EventTime = msg.EventTime
	}

	if len(msg.ReplicationClusters) > 0 {
		pulsarMsg.ReplicationClusters = msg.ReplicationClusters
	}

	if msg.DisableReplication {
		pulsarMsg.DisableReplication = msg.DisableReplication
	}

	messageID, err := producer.Send(cp.ctx, pulsarMsg)

	result := &ProducerResult{
		Message:     msg,
		MessageID:   messageID,
		Error:       err,
		PublishTime: time.Now(),
	}

	if err != nil {
		cp.client.logger.Errorf("Failed to send message to topic %s: %v", topic, err)
	} else {
		cp.client.logger.Debugf("Message sent successfully to topic %s, MessageID: %v", topic, messageID)
	}

	return result
}

// getOrCreateProducer gets an existing producer or creates a new one for the given topic
func (cp *ConcurrentProducer) getOrCreateProducer(topic string) (pulsar.Producer, error) {
	cp.mu.RLock()
	producer, exists := cp.producers[topic]
	cp.mu.RUnlock()

	if exists {
		return producer, nil
	}

	cp.mu.Lock()
	defer cp.mu.Unlock()

	// Double-check pattern
	if producer, exists := cp.producers[topic]; exists {
		return producer, nil
	}

	producerOptions := pulsar.ProducerOptions{
		Topic:                   topic,
		Name:                    cp.config.ProducerName,
		SendTimeout:             cp.config.SendTimeout,
		DisableBatching:         cp.config.DisableBatching,
		BatchingMaxPublishDelay: cp.config.BatchingMaxPublishDelay,
		BatchingMaxMessages:     cp.config.BatchingMaxMessages,
		BatchingMaxSize:         cp.config.BatchingMaxSize,
		MaxPendingMessages:      cp.config.MaxPendingMessages,
		HashingScheme:           cp.config.HashingScheme,
		CompressionType:         cp.config.CompressionType,
		CompressionLevel:        cp.config.CompressionLevel,
	}

	newProducer, err := cp.client.client.CreateProducer(producerOptions)
	if err != nil {
		return nil, err
	}

	cp.producers[topic] = newProducer
	cp.client.logger.Infof("Created new producer for topic: %s", topic)
	return newProducer, nil
}

// SendChannel returns the channel for sending messages
func (cp *ConcurrentProducer) SendChannel() chan<- *ProducerMessage {
	return cp.inputCh
}

// ResultChannel returns the channel for receiving send results
func (cp *ConcurrentProducer) ResultChannel() <-chan *ProducerResult {
	return cp.resultCh
}

// SendAsync sends a message asynchronously
func (cp *ConcurrentProducer) SendAsync(msg *ProducerMessage) error {
	select {
	case cp.inputCh <- msg:
		return nil
	case <-cp.ctx.Done():
		return fmt.Errorf("producer is closed")
	default:
		return fmt.Errorf("send queue is full")
	}
}

// SendSync sends a message synchronously and waits for the result
func (cp *ConcurrentProducer) SendSync(msg *ProducerMessage, timeout time.Duration) (*ProducerResult, error) {
	if err := cp.SendAsync(msg); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case result := <-cp.resultCh:
			if result.Message == msg {
				return result, nil
			}
			// This result belongs to another message, put it back (non-blocking)
			select {
			case cp.resultCh <- result:
			default:
				cp.client.logger.Warn("Could not return result to channel")
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for send result")
		case <-cp.ctx.Done():
			return nil, fmt.Errorf("producer is closed")
		}
	}
}

// Close closes the concurrent producer
func (cp *ConcurrentProducer) Close() error {
	cp.cancel()
	close(cp.inputCh)

	// Wait for workers to finish
	cp.wg.Wait()

	// Close all producers
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for topic, producer := range cp.producers {
		producer.Close()
		cp.client.logger.Infof("Closed producer for topic: %s", topic)
	}

	close(cp.resultCh)
	cp.client.logger.Info("Concurrent producer closed")
	return nil
}

// Stats returns basic statistics about the producer
func (cp *ConcurrentProducer) Stats() map[string]interface{} {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	return map[string]interface{}{
		"topics":       len(cp.producers),
		"workers":      cp.workerCount,
		"input_queue":  len(cp.inputCh),
		"result_queue": len(cp.resultCh),
	}
}
