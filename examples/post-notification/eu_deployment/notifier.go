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

func (r *notifier) Notify(ctx context.Context, postId string, userId int) error {

	message := Message{postId, userId}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Convert the byte slice to a string
	messageString := string(messageJSON)

	ctx, err = r.clientRabbitMQ.Write(ctx, "notifications", "", messageString)
	if err != nil {
		return err
	}

	return nil
}
