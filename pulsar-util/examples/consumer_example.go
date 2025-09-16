package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/drlucaa/dotfiles/pulsar-util"
)

func main() {
	// Create client
	client, err := pulsarutil.NewClientBuilder().
		ServiceURL("pulsar://localhost:6650").
		Build()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Define message handler
	messageHandler := func(ctx context.Context, msg *pulsarutil.ConsumerMessage) error {
		log.Printf("Received message from topic %s: %s (Key: %s, MessageID: %v)",
			msg.Topic, string(msg.Payload), msg.Key, msg.MessageID)

		// Simulate some processing time
		time.Sleep(100 * time.Millisecond)

		// Return nil for successful processing (message will be auto-acked)
		// Return error to nack the message
		return nil
	}

	// Create consumer using builder pattern with handler
	consumer, err := client.NewConsumerBuilder([]string{"my-topic"}, "my-subscription").
		SubscriptionType(pulsar.Shared).
		SubscriptionInitialPosition(pulsar.SubscriptionPositionEarliest).
		ReceiverQueueSize(1000).
		MaxUnackedMessages(1000).
		AckTimeout(30*time.Second).
		Build(messageHandler, true, 5) // 5 worker goroutines, auto-ack enabled
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Print consumer stats periodically
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				stats := consumer.Stats()
				log.Printf("Consumer stats: %+v", stats)
			}
		}
	}()

	// Setup graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Consumer started. Press Ctrl+C to exit.")

	// Wait for shutdown signal
	<-sigCh
	log.Println("Shutting down consumer...")
}
