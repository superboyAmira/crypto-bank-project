package service

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var (
	// Notification metrics
	NotificationsSent = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "notifications_sent_total",
			Help: "Total number of notifications sent",
		},
		[]string{"type", "channel"},
	)

	NotificationsFailed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "notifications_failed_total",
			Help: "Total number of failed notifications",
		},
		[]string{"type", "channel"},
	)
)

func init() {
	prometheus.MustRegister(NotificationsSent)
	prometheus.MustRegister(NotificationsFailed)
	
	// Initialize metrics with zero values to make them visible
	NotificationsSent.WithLabelValues("transaction", "email").Add(0)
	NotificationsSent.WithLabelValues("exchange", "push").Add(0)
	NotificationsSent.WithLabelValues("account", "email").Add(0)
	NotificationsSent.WithLabelValues("wallet", "email").Add(0)
}

type NotificationService struct {
	logger        *zap.Logger
	notifications []Notification
	mu            sync.RWMutex
}

type Notification struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Channel   string `json:"channel"`
	Timestamp int64  `json:"timestamp"`
}

func NewNotificationService(logger *zap.Logger) *NotificationService {
	return &NotificationService{
		logger:        logger,
		notifications: make([]Notification, 0),
	}
}

// ProcessTransactionEvent processes transaction events and sends notifications
func (s *NotificationService) ProcessTransactionEvent(body []byte) error {
	var event struct {
		TransactionID string  `json:"transaction_id"`
		UserID        string  `json:"user_id"`
		Type          string  `json:"type"`
		Amount        float64 `json:"amount"`
		Currency      string  `json:"currency"`
		Status        string  `json:"status"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		s.logger.Error("Failed to unmarshal transaction event", zap.Error(err))
		return err
	}

	// Create notification message
	var title, message string
	switch event.Type {
	case "TRANSFER":
		title = "Transfer Completed"
		message = fmt.Sprintf("Your transfer of %.2f %s has been %s", event.Amount, event.Currency, event.Status)
	case "DEPOSIT":
		title = "Deposit Received"
		message = fmt.Sprintf("Deposit of %.2f %s has been credited to your account", event.Amount, event.Currency)
	case "WITHDRAW":
		title = "Withdrawal Processed"
		message = fmt.Sprintf("Withdrawal of %.2f %s has been processed", event.Amount, event.Currency)
	case "EXCHANGE":
		title = "Exchange Completed"
		message = fmt.Sprintf("Exchange transaction of %.2f %s has been completed", event.Amount, event.Currency)
	default:
		title = "Transaction Update"
		message = fmt.Sprintf("Transaction of %.2f %s status: %s", event.Amount, event.Currency, event.Status)
	}

	// Send notifications via different channels
	s.sendNotification(event.UserID, "transaction", title, message, "email")
	s.sendNotification(event.UserID, "transaction", title, message, "push")

	return nil
}

// ProcessExchangeEvent processes exchange events and sends notifications
func (s *NotificationService) ProcessExchangeEvent(body []byte) error {
	var event struct {
		ExchangeID   string  `json:"exchange_id"`
		UserID       string  `json:"user_id"`
		FromCurrency string  `json:"from_currency"`
		ToCurrency   string  `json:"to_currency"`
		FromAmount   float64 `json:"from_amount"`
		ToAmount     float64 `json:"to_amount"`
		Status       string  `json:"status"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		s.logger.Error("Failed to unmarshal exchange event", zap.Error(err))
		return err
	}

	title := "Exchange Completed"
	message := fmt.Sprintf("Successfully exchanged %.6f %s to %.6f %s",
		event.FromAmount, event.FromCurrency, event.ToAmount, event.ToCurrency)

	s.sendNotification(event.UserID, "exchange", title, message, "email")
	s.sendNotification(event.UserID, "exchange", title, message, "push")

	return nil
}

// ProcessAccountEvent processes account creation events
func (s *NotificationService) ProcessAccountEvent(body []byte) error {
	var event struct {
		AccountID string `json:"account_id"`
		UserID    string `json:"user_id"`
		Currency  string `json:"currency"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		s.logger.Error("Failed to unmarshal account event", zap.Error(err))
		return err
	}

	title := "New Account Created"
	message := fmt.Sprintf("Your new %s account has been created successfully", event.Currency)

	s.sendNotification(event.UserID, "account", title, message, "email")

	return nil
}

// ProcessWalletEvent processes wallet creation events
func (s *NotificationService) ProcessWalletEvent(body []byte) error {
	var event struct {
		WalletID   string `json:"wallet_id"`
		UserID     string `json:"user_id"`
		CryptoType string `json:"crypto_type"`
	}

	if err := json.Unmarshal(body, &event); err != nil {
		s.logger.Error("Failed to unmarshal wallet event", zap.Error(err))
		return err
	}

	title := "New Crypto Wallet Created"
	message := fmt.Sprintf("Your new %s wallet has been created successfully", event.CryptoType)

	s.sendNotification(event.UserID, "wallet", title, message, "email")

	return nil
}

// sendNotification simulates sending a notification
func (s *NotificationService) sendNotification(userID, notificationType, title, message, channel string) {
	notification := Notification{
		ID:        fmt.Sprintf("%s-%d", notificationType, len(s.notifications)+1),
		Type:      notificationType,
		UserID:    userID,
		Title:     title,
		Message:   message,
		Channel:   channel,
		Timestamp: 0, // Should be time.Now().Unix() in production
	}

	s.mu.Lock()
	s.notifications = append(s.notifications, notification)
	s.mu.Unlock()

	// In production, this would send actual emails, push notifications, SMS, etc.
	s.logger.Info("Notification sent",
		zap.String("user_id", userID),
		zap.String("type", notificationType),
		zap.String("channel", channel),
		zap.String("title", title),
	)

	NotificationsSent.WithLabelValues(notificationType, channel).Inc()
}

// GetNotifications returns all notifications (for testing/debugging)
func (s *NotificationService) GetNotifications() []Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()

	notificationsCopy := make([]Notification, len(s.notifications))
	copy(notificationsCopy, s.notifications)

	return notificationsCopy
}

// GetUserNotifications returns notifications for a specific user
func (s *NotificationService) GetUserNotifications(userID string) []Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var userNotifications []Notification
	for _, notification := range s.notifications {
		if notification.UserID == userID {
			userNotifications = append(userNotifications, notification)
		}
	}

	return userNotifications
}

