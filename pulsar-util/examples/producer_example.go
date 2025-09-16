package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/drlucaa/dotfiles/pulsar-util"
)

func main() {
	// Create client using builder pattern
	client, err := pulsarutil.NewClientBuilder().
		ServiceURL("pulsar://localhost:6650").
		ConnectionTimeout(5 * time.Second).
		OperationTimeout(30 * time.Second).
		Build()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Create producer using builder pattern
	producer, err := client.NewProducerBuilder("my-topic").
		ProducerName("example-producer").
		SendTimeout(10 * time.Second).
		BatchingMaxPublishDelay(100 * time.Millisecond).
		CompressionType(pulsar.LZ4).
		Build(3) // 3 worker goroutines
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start result processor
	go func() {
		for {
			select {
			case result := <-producer.ResultChannel():
				if result.Error != nil {
					log.Printf("Failed to send message: %v", result.Error)
				} else {
					log.Printf("Message sent successfully: %v", result.MessageID)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Send messages
	sendCh := producer.SendChannel()

	// Send some test messages
	for i := 0; i < 10; i++ {
		msg := &pulsarutil.ProducerMessage{
			Payload: []byte(fmt.Sprintf("Hello World %d", i)),
			Key:     fmt.Sprintf("key-%d", i),
			Properties: map[string]string{
				"index": fmt.Sprintf("%d", i),
				"type":  "test",
			},
		}

		select {
		case sendCh <- msg:
			log.Printf("Queued message %d", i)
		case <-ctx.Done():
			return
		default:
			log.Printf("Send queue full, skipping message %d", i)
		}

		time.Sleep(500 * time.Millisecond)
	}

	// Send a synchronous message
	syncMsg := &pulsarutil.ProducerMessage{
		Payload: []byte("Synchronous message"),
		Key:     "sync-key",
	}

	result, err := producer.SendSync(syncMsg, 5*time.Second)
	if err != nil {
		log.Printf("Sync send failed: %v", err)
	} else {
		log.Printf("Sync message sent: %v", result.MessageID)
	}

	// Print producer stats
	stats := producer.Stats()
	log.Printf("Producer stats: %+v", stats)

	// Wait for shutdown signal
	<-sigCh
	log.Println("Shutting down...")
}
