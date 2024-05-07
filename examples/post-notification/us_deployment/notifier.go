package main

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/TiagoMalhadas/xcweaver"
)

type Message struct {
	PostId                  string
	UserId                  int
	StartTimeMs             int64
	NotificationStartTimeMs int64
}

// Notifier component.
type Notifier interface {
}

type NotifierOptions struct {
	NumWorkers int `toml:"num_workers"`
}

// Implementation of the Notifier component.
type notifier struct {
	xcweaver.Implements[Notifier]
	xcweaver.WithConfig[NotifierOptions]
	follower_Notify xcweaver.Ref[Follower_Notify]
	clientRabbitMQ  xcweaver.Antipode
}

func (n *notifier) Init(ctx context.Context) error {
	logger := n.Logger(ctx)
	logger.Info("notifier service at us running!")

	var wg sync.WaitGroup
	wg.Add(n.Config().NumWorkers)
	for i := 1; i <= n.Config().NumWorkers; i++ {
		go func(i int) {
			defer wg.Done()
			err := n.workerThread(ctx, i)
			logger.Error("error in worker thread", "msg", err.Error())
		}(i)
	}
	wg.Wait()
	return nil
}

func (n *notifier) workerThread(ctx context.Context, workerid int) error {
	logger := n.Logger(ctx)

	var forever chan struct{}

	go func() {
		for {
			notification, lineage, err := n.clientRabbitMQ.Read(ctx, "notifications", "")
			if err != nil {
				logger.Error("Failed to read from rabbitMQ!", "msg", err.Error(), "workerId", workerid)
				continue
			}

			var decodedNotification Message
			err = json.Unmarshal([]byte(notification), &decodedNotification)
			if err != nil {
				logger.Error("Error unmarshaling JSON!", "msg", err.Error(), "workerId", workerid)
				continue
			}

			postId := decodedNotification.PostId
			userId := decodedNotification.UserId
			startTimeMs := decodedNotification.StartTimeMs
			notificationStartTimeMs := decodedNotification.NotificationStartTimeMs
			logger.Debug("New notification received", "postId", postId, "userId", userId, "workerId", workerid)
			queueDurationMs.Put(float64(time.Now().UnixMilli() - notificationStartTimeMs))
			notificationsReceived.Inc()

			ctx, err = xcweaver.Transfer(ctx, lineage)
			if err != nil {
				logger.Error("Error transfering the new lineage to context!", "msg", err.Error(), "workerId", workerid)
				continue
			}

			err = n.follower_Notify.Get().Follower_Notify(ctx, postId, userId, startTimeMs)
			if err != nil {
				continue
			}
		}
	}()

	<-forever
	return nil
}
