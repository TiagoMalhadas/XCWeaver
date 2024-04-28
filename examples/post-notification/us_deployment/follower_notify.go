package main

import (
	"context"
	"time"

	"github.com/TiagoMalhadas/xcweaver"
)

// Follower_Notify component.
type Follower_Notify interface {
	Follower_Notify(context.Context, string, int, int64) error
}

// Implementation of the Follower_Notify component.
type follower_Notify struct {
	xcweaver.Implements[Follower_Notify]
	post_storage xcweaver.Ref[Post_storage_america]
}

func (f *follower_Notify) Init(ctx context.Context) error {
	logger := f.Logger(ctx)

	logger.Info("follower notify service running!")
	return nil
}

func (f *follower_Notify) Follower_Notify(ctx context.Context, postId string, userId int, startTimeMs int64) error {
	logger := f.Logger(ctx)

	post, err := f.post_storage.Get().GetPost(ctx, postId)
	postNotificationDuration.Put(float64(time.Now().UnixMilli() - startTimeMs))
	if err != nil {
		return err
	}

	logger.Debug("Post read successfully!", "post", post)

	return nil
}
