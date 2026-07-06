package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageQueue implements message queue for distributed processing
type MessageQueue struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// Queue names
const (
	QueueBets          = "tigercasino.bets"
	QueueWithdrawals   = "tigercasino.withdrawals"
	QueueDeposits      = "tigercasino.deposits"
	QueueGameEvents    = "tigercasino.game_events"
	QueueNotifications = "tigercasino.notifications"
	QueueAudit         = "tigercasino.audit"
)

// NewMessageQueue creates a new RabbitMQ connection
func NewMessageQueue(url string) (*MessageQueue, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	mq := &MessageQueue{
		conn:    conn,
		channel: ch,
	}

	// Declare queues
	queues := []string{QueueBets, QueueWithdrawals, QueueDeposits, QueueGameEvents, QueueNotifications, QueueAudit}
	for _, q := range queues {
		if err := ch.QueueDeclare(q, true, false, false, false, nil); err != nil {
			mq.Close()
			return nil, fmt.Errorf("failed to declare queue %s: %w", q, err)
		}
	}

	return mq, nil
}

// Publish publishes a message to a queue
func (m *MessageQueue) Publish(ctx context.Context, queue string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return m.channel.PublishWithContext(ctx, "", queue, false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		})
}

// Consume starts consuming messages from a queue
func (m *MessageQueue) Consume(queue string, handler func([]byte) error) error {
	msgs, err := m.channel.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				msg.Nack(false, true) // Requeue on error
			} else {
				msg.Ack(false)
			}
		}
	}()

	return nil
}

// Close closes the connection
func (m *MessageQueue) Close() error {
	if m.channel != nil {
		m.channel.Close()
	}
	if m.conn != nil {
		return m.conn.Close()
	}
	return nil
}

// BetMessage represents a bet transaction message
type BetMessage struct {
	UserID     string    `json:"user_id"`
	GameID     string    `json:"game_id"`
	RoundID    string    `json:"round_id"`
	Amount     float64   `json:"amount"`
	Multiplier float64   `json:"multiplier,omitempty"`
	WinAmount  float64   `json:"win_amount,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// WithdrawalMessage represents a withdrawal request
type WithdrawalMessage struct {
	UserID        string    `json:"user_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Network       string    `json:"network"`
	Address       string    `json:"address"`
	TransactionID string    `json:"transaction_id,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
}

// DepositMessage represents a deposit notification
type DepositMessage struct {
	UserID           string    `json:"user_id"`
	Amount           float64   `json:"amount"`
	Currency         string    `json:"currency"`
	Network          string    `json:"network"`
	TransactionHash string    `json:"transaction_hash"`
	Confirmations    int       `json:"confirmations"`
	Timestamp        time.Time `json:"timestamp"`
}

// GameEvent represents a game event for analytics
type GameEvent struct {
	EventType string          `json:"event_type"`
	UserID    string          `json:"user_id"`
	GameID    string          `json:"game_id"`
	RoundID   string          `json:"round_id"`
	Data      json.RawMessage `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
}
