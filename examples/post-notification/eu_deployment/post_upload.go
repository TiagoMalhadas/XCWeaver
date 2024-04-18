package main

import (
	"context"

	"github.com/TiagoMalhadas/xcweaver"
)

// Post_upload component.
type Post_upload interface {
	Post(context.Context, string, int) error
}

// Implementation of the Post_upload component.
type post_upload struct {
	xcweaver.Implements[Post_upload]
	post_storage xcweaver.Ref[Post_storage_europe]
	notifier     xcweaver.Ref[Notifier]
}

func (p *post_upload) Init(ctx context.Context) error {
	logger := p.Logger(ctx)
	logger.Info("post upload service at eu running!")

	return nil
}

// Forwards the post data to Post_storage component and then sends the post id to
// the Notifier component
func (p *post_upload) Post(ctx context.Context, post string, userId int) error {
	logger := p.Logger(ctx)

	//send post to post_storage
	lineage, postId, err := p.post_storage.Get().Post(ctx, post)
	if err != nil {
		return err
	}

	ctx, err = xcweaver.Transfer(ctx, lineage)
	if err != nil {
		logger.Error("Error transfering lineage to context!", "msg", err.Error())
		return err
	}

	//send postID and userId to notifier
	err = p.notifier.Get().Notify(ctx, postId, userId)
	if err != nil {
		return err
	}

	return nil
}
