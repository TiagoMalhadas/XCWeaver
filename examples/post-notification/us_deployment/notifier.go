package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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

func (r *notifier) Init(ctx context.Context) error {

	var forever chan struct{}

	go func() {
		for {
			notification, lineage, err := r.clientRabbitMQ.Read(ctx, "notifications", "")
			if err != nil {
				log.Printf("Failed to read from rabbitMQ: %v", err)
				continue
			}

			var decodedNotification Message
			err = json.Unmarshal([]byte(notification), &decodedNotification)
			if err != nil {
				log.Printf("Error unmarshaling JSON: %v", err)
				continue
			}

			postId := decodedNotification.PostId
			userId := decodedNotification.UserId

			ctx, err = xcweaver.Transfer(ctx, lineage)

			if err != nil {
				fmt.Println(err)
				continue
			}

			err = r.follower_Notify.Get().Follower_Notify(ctx, postId, userId)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}()

	<-forever
	return nil
}
