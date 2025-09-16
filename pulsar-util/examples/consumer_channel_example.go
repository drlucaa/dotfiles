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

	// Create channel-based consumer (no handler, using channels directly)
	consumer, err := client.NewConsumerBuilder([]string{"my-topic"}, "channel-subscription").
		SubscriptionType(pulsar.Shared).
		SubscriptionInitialPosition(pulsar.SubscriptionPositionLatest).
		ReceiverQueueSize(1000).
		BuildChannelBased(0) // No workers since we're using channels directly
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Process messages from channel
	go func() {
		messageCh := consumer.MessageChannel()

		for {
			select {
			case msg, ok := <-messageCh:
				if !ok {
					log.Println("Message channel closed")
					return
				}

				log.Printf("Processing message from topic %s: %s (Key: %s)",
					msg.Topic, string(msg.Payload), msg.Key)

				// Simulate processing
				time.Sleep(100 * time.Millisecond)

				// Manually acknowledge the message
				if err := consumer.Ack(msg); err != nil {
					log.Printf("Failed to ack message: %v", err)
				} else {
					log.Printf("Acknowledged message: %v", msg.MessageID)
				}

			case <-ctx.Done():
				log.Println("Message processor stopping")
				return
			}
		}
	}()

	log.Println("Channel-based consumer started. Press Ctrl+C to exit.")

	// Wait for shutdown signal
	<-sigCh
	log.Println("Shutting down...")
	cancel()
}
