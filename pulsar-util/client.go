package pulsarutil

import (
	"fmt"
	"log"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

// Client wraps the Apache Pulsar client with additional functionality
type Client struct {
	client   pulsar.Client
	config   *Config
	logger   Logger
	closeCh  chan struct{}
	isClosed bool
}

// Config holds configuration for the Pulsar client
type Config struct {
	ServiceURL                 string
	ConnectionTimeout          time.Duration
	OperationTimeout           time.Duration
	MaxConnectionsPerBroker    int
	Authentication             pulsar.Authentication
	TLSAllowInsecureConnection bool
	TLSTrustCertsFilePath      string
}

// Logger defines the logging interface
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// DefaultLogger provides a basic logger implementation
type DefaultLogger struct{}

func (l *DefaultLogger) Debug(args ...interface{}) {
	log.Println(append([]interface{}{"[DEBUG]"}, args...)...)
}
func (l *DefaultLogger) Info(args ...interface{}) {
	log.Println(append([]interface{}{"[INFO]"}, args...)...)
}
func (l *DefaultLogger) Warn(args ...interface{}) {
	log.Println(append([]interface{}{"[WARN]"}, args...)...)
}
func (l *DefaultLogger) Error(args ...interface{}) {
	log.Println(append([]interface{}{"[ERROR]"}, args...)...)
}
func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	log.Printf("[DEBUG] "+format, args...)
}
func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}
func (l *DefaultLogger) Warnf(format string, args ...interface{}) {
	log.Printf("[WARN] "+format, args...)
}
func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		ServiceURL:                 "pulsar://localhost:6650",
		ConnectionTimeout:          5 * time.Second,
		OperationTimeout:           30 * time.Second,
		MaxConnectionsPerBroker:    1,
		TLSAllowInsecureConnection: false,
	}
}

// NewClient creates a new Pulsar client with the given configuration
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	clientOptions := pulsar.ClientOptions{
		URL:                        config.ServiceURL,
		ConnectionTimeout:          config.ConnectionTimeout,
		OperationTimeout:           config.OperationTimeout,
		MaxConnectionsPerBroker:    config.MaxConnectionsPerBroker,
		TLSAllowInsecureConnection: config.TLSAllowInsecureConnection,
	}

	if config.Authentication != nil {
		clientOptions.Authentication = config.Authentication
	}

	if config.TLSTrustCertsFilePath != "" {
		clientOptions.TLSTrustCertsFilePath = config.TLSTrustCertsFilePath
	}

	pulsarClient, err := pulsar.NewClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create Pulsar client: %w", err)
	}

	client := &Client{
		client:  pulsarClient,
		config:  config,
		logger:  &DefaultLogger{},
		closeCh: make(chan struct{}),
	}

	return client, nil
}

// SetLogger sets a custom logger for the client
func (c *Client) SetLogger(logger Logger) {
	c.logger = logger
}

// Close closes the Pulsar client and all associated resources
func (c *Client) Close() error {
	if c.isClosed {
		return nil
	}

	c.isClosed = true
	close(c.closeCh)

	if c.client != nil {
		c.client.Close()
	}

	c.logger.Info("Pulsar client closed")
	return nil
}

// IsClosed returns true if the client is closed
func (c *Client) IsClosed() bool {
	return c.isClosed
}

// GetNativeClient returns the underlying Pulsar client for advanced use cases
func (c *Client) GetNativeClient() pulsar.Client {
	return c.client
}
