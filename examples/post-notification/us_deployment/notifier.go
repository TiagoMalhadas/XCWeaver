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

// Notifier component.
type Notifier interface {
}

// Implementation of the Notifier component.
type notifier struct {
	xcweaver.Implements[Notifier]
	follower_Notify xcweaver.Ref[Follower_Notify]
	clientRabbitMQ  xcweaver.Antipode
}

func (n *notifier) Init(ctx context.Context) error {
	logger := n.Logger(ctx)
	logger.Info("notifier service at us running!")

	var forever chan struct{}

	go func() {
		for {
			notification, lineage, err := n.clientRabbitMQ.Read(ctx, "notifications", "")
			if err != nil {
				logger.Error("Failed to read from rabbitMQ!", "msg", err.Error())
				continue
			}

			var decodedNotification Message
			err = json.Unmarshal([]byte(notification), &decodedNotification)
			if err != nil {
				logger.Error("Error unmarshaling JSON!", "msg", err.Error())
				continue
			}

			postId := decodedNotification.PostId
			userId := decodedNotification.UserId
			logger.Debug("New notification received", "postId", postId, "userId", userId)
			notificationsReceived.Inc()

			ctx, err = xcweaver.Transfer(ctx, lineage)
			if err != nil {
				logger.Error("Error transfering the new lineage to context!", "msg", err.Error())
				continue
			}

			err = n.follower_Notify.Get().Follower_Notify(ctx, postId, userId)
			if err != nil {
				continue
			}
		}
	}()

	<-forever
	return nil
}
