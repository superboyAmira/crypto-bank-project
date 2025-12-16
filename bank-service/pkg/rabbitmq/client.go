package rabbitmq

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/crypto-bank/bank-service/pkg/logger"
	"go.uber.org/zap"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewClient creates a new RabbitMQ client
func NewClient(url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Client{
		conn:    conn,
		channel: channel,
	}, nil
}

// DeclareExchange declares an exchange
func (c *Client) DeclareExchange(name, kind string) error {
	return c.channel.ExchangeDeclare(
		name,  // name
		kind,  // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
}

// DeclareQueue declares a queue
func (c *Client) DeclareQueue(name string) (amqp.Queue, error) {
	return c.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

// BindQueue binds a queue to an exchange
func (c *Client) BindQueue(queueName, exchangeName, routingKey string) error {
	return c.channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)
}

// PublishEvent publishes an event to an exchange
func (c *Client) PublishEvent(exchange, routingKey string, event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = c.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	logger.Debug("Event published",
		zap.String("exchange", exchange),
		zap.String("routing_key", routingKey),
	)

	return nil
}

// Close closes the connection and channel
func (c *Client) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}

// Event types
const (
	ExchangeEvents = "bank.events"
	
	// Routing keys
	EventTransactionCreated = "transaction.created"
	EventTransactionCompleted = "transaction.completed"
	EventExchangeCreated = "exchange.created"
	EventExchangeCompleted = "exchange.completed"
	EventAccountCreated = "account.created"
	EventWalletCreated = "wallet.created"
)

// Event structures
type TransactionEvent struct {
	TransactionID string  `json:"transaction_id"`
	UserID        string  `json:"user_id"`
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
}

type ExchangeEvent struct {
	ExchangeID   string  `json:"exchange_id"`
	UserID       string  `json:"user_id"`
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	FromAmount   float64 `json:"from_amount"`
	ToAmount     float64 `json:"to_amount"`
	Status       string  `json:"status"`
}

type AccountEvent struct {
	AccountID string `json:"account_id"`
	UserID    string `json:"user_id"`
	Currency  string `json:"currency"`
}

type WalletEvent struct {
	WalletID   string `json:"wallet_id"`
	UserID     string `json:"user_id"`
	CryptoType string `json:"crypto_type"`
}

