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

// Forwards the post data to Post_storage component and then sends the post id to
// the Notifier component
func (r *post_upload) Post(ctx context.Context, post string, userId int) error {

	//send post to post_storage
	var post_storage Post_storage_europe = r.post_storage.Get()
	lineage, postId, err := post_storage.Post(ctx, post)
	if err != nil {
		return err
	}

	ctx, err = xcweaver.Transfer(ctx, lineage)
	if err != nil {
		return err
	}

	//send postID and userId to notifier
	var notifier Notifier = r.notifier.Get()
	err = notifier.Notify(ctx, postId, userId)
	if err != nil {
		return err
	}

	return nil
}
