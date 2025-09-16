package pulsarutil

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.ServiceURL != "pulsar://localhost:6650" {
		t.Errorf("Expected ServiceURL to be 'pulsar://localhost:6650', got %s", config.ServiceURL)
	}

	if config.ConnectionTimeout != 5*time.Second {
		t.Errorf("Expected ConnectionTimeout to be 5s, got %v", config.ConnectionTimeout)
	}

	if config.OperationTimeout != 30*time.Second {
		t.Errorf("Expected OperationTimeout to be 30s, got %v", config.OperationTimeout)
	}
}

func TestClientBuilder(t *testing.T) {
	builder := NewClientBuilder()

	// Test builder methods
	builder = builder.ServiceURL("pulsar://test:6650").
		ConnectionTimeout(10 * time.Second).
		OperationTimeout(60 * time.Second).
		MaxConnectionsPerBroker(5).
		TLSAllowInsecureConnection(true).
		TLSTrustCertsFilePath("/test/path")

	if builder.config.ServiceURL != "pulsar://test:6650" {
		t.Errorf("Expected ServiceURL to be 'pulsar://test:6650', got %s", builder.config.ServiceURL)
	}

	if builder.config.ConnectionTimeout != 10*time.Second {
		t.Errorf("Expected ConnectionTimeout to be 10s, got %v", builder.config.ConnectionTimeout)
	}

	if builder.config.OperationTimeout != 60*time.Second {
		t.Errorf("Expected OperationTimeout to be 60s, got %v", builder.config.OperationTimeout)
	}

	if builder.config.MaxConnectionsPerBroker != 5 {
		t.Errorf("Expected MaxConnectionsPerBroker to be 5, got %d", builder.config.MaxConnectionsPerBroker)
	}

	if !builder.config.TLSAllowInsecureConnection {
		t.Error("Expected TLSAllowInsecureConnection to be true")
	}

	if builder.config.TLSTrustCertsFilePath != "/test/path" {
		t.Errorf("Expected TLSTrustCertsFilePath to be '/test/path', got %s", builder.config.TLSTrustCertsFilePath)
	}
}

func TestDefaultProducerConfig(t *testing.T) {
	topic := "test-topic"
	config := DefaultProducerConfig(topic)

	if config.Topic != topic {
		t.Errorf("Expected Topic to be '%s', got %s", topic, config.Topic)
	}

	if config.SendTimeout != 30*time.Second {
		t.Errorf("Expected SendTimeout to be 30s, got %v", config.SendTimeout)
	}

	if config.DisableBatching {
		t.Error("Expected DisableBatching to be false")
	}

	if config.BatchingMaxPublishDelay != 10*time.Millisecond {
		t.Errorf("Expected BatchingMaxPublishDelay to be 10ms, got %v", config.BatchingMaxPublishDelay)
	}

	if config.BatchingMaxMessages != 1000 {
		t.Errorf("Expected BatchingMaxMessages to be 1000, got %d", config.BatchingMaxMessages)
	}

	if config.MaxPendingMessages != 1000 {
		t.Errorf("Expected MaxPendingMessages to be 1000, got %d", config.MaxPendingMessages)
	}
}

func TestDefaultConsumerConfig(t *testing.T) {
	topics := []string{"topic1", "topic2"}
	subscription := "test-subscription"
	config := DefaultConsumerConfig(topics, subscription)

	if len(config.Topics) != 2 || config.Topics[0] != "topic1" || config.Topics[1] != "topic2" {
		t.Errorf("Expected Topics to be %v, got %v", topics, config.Topics)
	}

	if config.Subscription != subscription {
		t.Errorf("Expected Subscription to be '%s', got %s", subscription, config.Subscription)
	}

	if config.MessageChannel != 1000 {
		t.Errorf("Expected MessageChannel to be 1000, got %d", config.MessageChannel)
	}

	if config.ReceiverQueueSize != 1000 {
		t.Errorf("Expected ReceiverQueueSize to be 1000, got %d", config.ReceiverQueueSize)
	}

	if config.MaxUnackedMessages != 1000 {
		t.Errorf("Expected MaxUnackedMessages to be 1000, got %d", config.MaxUnackedMessages)
	}

	if config.AckTimeout != 30*time.Second {
		t.Errorf("Expected AckTimeout to be 30s, got %v", config.AckTimeout)
	}

	if !config.RetryEnable {
		t.Error("Expected RetryEnable to be true")
	}
}

func TestDefaultLogger(t *testing.T) {
	logger := &DefaultLogger{}

	// Test that logger methods don't panic
	logger.Debug("test debug")
	logger.Info("test info")
	logger.Warn("test warn")
	logger.Error("test error")
	logger.Debugf("test debug %s", "formatted")
	logger.Infof("test info %s", "formatted")
	logger.Warnf("test warn %s", "formatted")
	logger.Errorf("test error %s", "formatted")
}

func TestProducerMessage(t *testing.T) {
	msg := &ProducerMessage{
		Topic:   "test-topic",
		Payload: []byte("test payload"),
		Key:     "test-key",
		Properties: map[string]string{
			"prop1": "value1",
			"prop2": "value2",
		},
		EventTime:           time.Now(),
		ReplicationClusters: []string{"cluster1", "cluster2"},
		DisableReplication:  true,
	}

	if msg.Topic != "test-topic" {
		t.Errorf("Expected Topic to be 'test-topic', got %s", msg.Topic)
	}

	if string(msg.Payload) != "test payload" {
		t.Errorf("Expected Payload to be 'test payload', got %s", string(msg.Payload))
	}

	if msg.Key != "test-key" {
		t.Errorf("Expected Key to be 'test-key', got %s", msg.Key)
	}

	if len(msg.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(msg.Properties))
	}

	if msg.Properties["prop1"] != "value1" {
		t.Errorf("Expected prop1 to be 'value1', got %s", msg.Properties["prop1"])
	}

	if len(msg.ReplicationClusters) != 2 {
		t.Errorf("Expected 2 replication clusters, got %d", len(msg.ReplicationClusters))
	}

	if !msg.DisableReplication {
		t.Error("Expected DisableReplication to be true")
	}
}

func TestConsumerMessage(t *testing.T) {
	now := time.Now()
	msg := &ConsumerMessage{
		Topic:           "test-topic",
		Subscription:    "test-subscription",
		Payload:         []byte("test payload"),
		Key:             "test-key",
		Properties:      map[string]string{"prop1": "value1"},
		PublishTime:     now,
		EventTime:       now,
		RedeliveryCount: 2,
	}

	if msg.Topic != "test-topic" {
		t.Errorf("Expected Topic to be 'test-topic', got %s", msg.Topic)
	}

	if msg.Subscription != "test-subscription" {
		t.Errorf("Expected Subscription to be 'test-subscription', got %s", msg.Subscription)
	}

	if string(msg.Payload) != "test payload" {
		t.Errorf("Expected Payload to be 'test payload', got %s", string(msg.Payload))
	}

	if msg.Key != "test-key" {
		t.Errorf("Expected Key to be 'test-key', got %s", msg.Key)
	}

	if msg.RedeliveryCount != 2 {
		t.Errorf("Expected RedeliveryCount to be 2, got %d", msg.RedeliveryCount)
	}
}
