# Pulsar Utility Library

A concurrent Apache Pulsar Go utility library that wraps the official Apache Pulsar Go client, providing a simplified, channel-based API for high-performance, concurrent message publishing and consuming.

## Features

- **Simplified API**: User-friendly interface that abstracts the complexity of the official Pulsar client
- **Concurrent Processing**: Safe, concurrent publishing and consuming using Go channels for maximum throughput
- **Channel-based Interface**: Unified channel-based API for both producers and consumers
- **Connection Management**: Automatic connection pooling and graceful error handling
- **Builder Pattern**: Fluent API for easy configuration
- **Flexible Message Handling**: Support for both callback-based and channel-based message processing
- **Comprehensive Monitoring**: Built-in statistics and logging
- **Graceful Shutdown**: Proper resource cleanup and graceful shutdown handling

## Installation

```bash
go get github.com/drlucaa/dotfiles/pulsar-util
```

## Quick Start

### Basic Producer

```go
package main

import (
    "log"
    "time"
    
    "github.com/drlucaa/dotfiles/pulsar-util"
)

func main() {
    // Create client
    client, err := pulsarutil.NewClientBuilder().
        ServiceURL("pulsar://localhost:6650").
        Build()
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Create producer
    producer, err := client.NewProducerBuilder("my-topic").
        ProducerName("my-producer").
        Build(3) // 3 worker goroutines
    if err != nil {
        log.Fatal(err)
    }
    defer producer.Close()

    // Send message
    msg := &pulsarutil.ProducerMessage{
        Payload: []byte("Hello World"),
        Key:     "my-key",
    }
    
    if err := producer.SendAsync(msg); err != nil {
        log.Printf("Failed to send: %v", err)
    }
}
```

### Basic Consumer

```go
package main

import (
    "context"
    "log"
    
    "github.com/apache/pulsar-client-go/pulsar"
    "github.com/drlucaa/dotfiles/pulsar-util"
)

func main() {
    client, err := pulsarutil.NewClientBuilder().
        ServiceURL("pulsar://localhost:6650").
        Build()
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Define message handler
    handler := func(ctx context.Context, msg *pulsarutil.ConsumerMessage) error {
        log.Printf("Received: %s", string(msg.Payload))
        return nil // nil = success, error = nack
    }

    // Create consumer
    consumer, err := client.NewConsumerBuilder([]string{"my-topic"}, "my-subscription").
        SubscriptionType(pulsar.Shared).
        Build(handler, true, 5) // auto-ack enabled, 5 workers
    if err != nil {
        log.Fatal(err)
    }
    defer consumer.Close()

    // Consumer will run until Close() is called
    select {} // Block forever
}
```

## API Reference

### Client

The `Client` wraps the Apache Pulsar client with additional functionality:

```go
// Create client with builder pattern
client, err := pulsarutil.NewClientBuilder().
    ServiceURL("pulsar://localhost:6650").
    ConnectionTimeout(5 * time.Second).
    OperationTimeout(30 * time.Second).
    MaxConnectionsPerBroker(10).
    TLSAllowInsecureConnection(false).
    Build()
```

### Producer

The `ConcurrentProducer` provides thread-safe, high-performance message publishing:

#### Configuration Options

```go
producer, err := client.NewProducerBuilder("my-topic").
    ProducerName("my-producer").
    SendTimeout(10 * time.Second).
    BatchingMaxPublishDelay(100 * time.Millisecond).
    BatchingMaxMessages(1000).
    CompressionType(pulsar.LZ4).
    MaxPendingMessages(1000).
    Build(workerCount)
```

#### Sending Messages

**Asynchronous (Recommended)**:
```go
msg := &pulsarutil.ProducerMessage{
    Payload: []byte("Hello World"),
    Key:     "message-key",
    Properties: map[string]string{
        "type": "greeting",
    },
}

// Non-blocking send
err := producer.SendAsync(msg)

// Process results
go func() {
    for result := range producer.ResultChannel() {
        if result.Error != nil {
            log.Printf("Send failed: %v", result.Error)
        } else {
            log.Printf("Sent: %v", result.MessageID)
        }
    }
}()
```

**Synchronous**:
```go
result, err := producer.SendSync(msg, 5*time.Second)
if err != nil {
    log.Printf("Send failed: %v", err)
} else {
    log.Printf("Sent: %v", result.MessageID)
}
```

**Channel-based**:
```go
sendCh := producer.SendChannel()
sendCh <- msg
```

### Consumer

The `ConcurrentConsumer` provides concurrent message consumption with two modes:

#### Handler-based (Recommended)

```go
handler := func(ctx context.Context, msg *pulsarutil.ConsumerMessage) error {
    // Process message
    log.Printf("Received: %s", string(msg.Payload))
    
    // Return nil for success (auto-ack if enabled)
    // Return error to nack the message
    return nil
}

consumer, err := client.NewConsumerBuilder([]string{"topic1", "topic2"}, "my-subscription").
    SubscriptionType(pulsar.Shared).
    SubscriptionInitialPosition(pulsar.SubscriptionPositionEarliest).
    ReceiverQueueSize(1000).
    MaxUnackedMessages(1000).
    AckTimeout(30 * time.Second).
    Build(handler, true, 5) // auto-ack, 5 workers
```

#### Channel-based

```go
consumer, err := client.NewConsumerBuilder([]string{"my-topic"}, "my-subscription").
    BuildChannelBased(0) // No workers, use channels directly

go func() {
    for msg := range consumer.MessageChannel() {
        // Process message
        log.Printf("Received: %s", string(msg.Payload))
        
        // Manual acknowledgment
        if err := consumer.Ack(msg); err != nil {
            log.Printf("Ack failed: %v", err)
        }
    }
}()
```

### Configuration

#### Client Configuration

```go
config := &pulsarutil.Config{
    ServiceURL:       "pulsar://localhost:6650",
    ConnectionTimeout: 5 * time.Second,
    OperationTimeout:  30 * time.Second,
    MaxConnectionsPerBroker: 10,
    TLSAllowInsecureConnection: false,
    TLSTrustCertsFilePath: "/path/to/certs",
    Authentication: nil, // pulsar.Authentication
}
```

#### Producer Configuration

```go
config := &pulsarutil.ProducerConfig{
    Topic:                   "my-topic",
    ProducerName:           "my-producer",
    SendTimeout:            30 * time.Second,
    DisableBatching:        false,
    BatchingMaxPublishDelay: 10 * time.Millisecond,
    BatchingMaxMessages:    1000,
    BatchingMaxSize:        128 * 1024, // 128KB
    MaxPendingMessages:     1000,
    HashingScheme:          pulsar.JavaStringHash,
    CompressionType:        pulsar.LZ4,
    CompressionLevel:       pulsar.Default,
}
```

#### Consumer Configuration

```go
config := &pulsarutil.ConsumerConfig{
    Topics:                     []string{"topic1", "topic2"},
    Subscription:               "my-subscription",
    SubscriptionType:           pulsar.Shared,
    SubscriptionInitialPosition: pulsar.SubscriptionPositionLatest,
    MessageChannel:             1000,
    ReceiverQueueSize:          1000,
    MaxUnackedMessages:         1000,
    AckTimeout:                 30 * time.Second,
    NackRedeliveryDelay:        1 * time.Minute,
    RetryEnable:                true,
}
```

## Advanced Usage

### Error Handling and Retry

The library automatically handles connection failures and provides retry mechanisms:

```go
// Configure DLQ (Dead Letter Queue) for failed messages
dlqPolicy := &pulsar.DLQPolicy{
    MaxDeliveries:   3,
    DeadLetterTopic: "my-topic-dlq",
}

consumer, err := client.NewConsumerBuilder([]string{"my-topic"}, "my-subscription").
    DLQ(dlqPolicy).
    RetryEnable(true).
    Build(handler, true, 5)
```

### Monitoring and Statistics

```go
// Producer stats
stats := producer.Stats()
log.Printf("Producer stats: %+v", stats)

// Consumer stats
stats = consumer.Stats()
log.Printf("Consumer stats: %+v", stats)
```

### Custom Logging

```go
type MyLogger struct{}

func (l *MyLogger) Info(args ...interface{}) {
    // Custom logging implementation
}
// ... implement other methods

client.SetLogger(&MyLogger{})
```

### Graceful Shutdown

```go
// Setup signal handling
sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

// Wait for signal
<-sigCh

// Graceful shutdown
producer.Close()
consumer.Close()
client.Close()
```

## Performance Tips

1. **Use appropriate worker counts**: Start with 2-5 workers per CPU core
2. **Configure batching**: Enable batching for high-throughput scenarios
3. **Tune queue sizes**: Adjust based on your memory constraints and latency requirements
4. **Use compression**: Enable LZ4 or ZSTD compression for large messages
5. **Monitor stats**: Use built-in statistics to tune performance

## Examples

See the `examples/` directory for complete working examples:

- `producer_example.go` - Complete producer example
- `consumer_example.go` - Handler-based consumer example  
- `consumer_channel_example.go` - Channel-based consumer example

## License

This library is provided as part of the dotfiles repository and follows the same licensing terms.