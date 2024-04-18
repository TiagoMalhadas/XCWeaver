package main

import (
	"context"
	"encoding/json"

	"github.com/TiagoMalhadas/xcweaver"
)

type Message struct {
	PostId string
	UserId int
}

// Server component.
type Notifier interface {
	Notify(context.Context, string, int) error
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

func (n *notifier) Notify(ctx context.Context, postId string, userId int) error {
	logger := n.Logger(ctx)

	message := Message{postId, userId}

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
