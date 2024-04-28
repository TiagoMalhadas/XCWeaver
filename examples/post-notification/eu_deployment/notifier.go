package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/TiagoMalhadas/xcweaver"
)

type Message struct {
	PostId                  string
	UserId                  int
	StartTimeMs             int64
	NotificationStartTimeMs int64
}

// Server component.
type Notifier interface {
	Notify(context.Context, string, int, int64) error
}

// Implementation of the Notifier component.
type notifier struct {
	xcweaver.Implements[Notifier]
	clientRabbitMQ xcweaver.Antipode
}

func (n *notifier) Init(ctx context.Context) error {
	logger := n.Logger(ctx)
	logger.Info("notifier service at eu running!")

	return nil
}

func (n *notifier) Notify(ctx context.Context, postId string, userId int, startTimeMs int64) error {
	logger := n.Logger(ctx)

	notificationStartTimeMs := time.Now().UnixMilli()
	message := Message{postId, userId, startTimeMs, notificationStartTimeMs}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Convert the byte slice to a string
	messageString := string(messageJSON)

	_, err = n.clientRabbitMQ.Write(ctx, "notifications", "", messageString)
	if err != nil {
		logger.Error("Error writing notification to queue", "msg", err.Error())
		return err
	}
	notificationsSent.Inc()
	logger.Debug("nofitication successfully written", "postId", postId, "userId", userId)

	return nil
}
