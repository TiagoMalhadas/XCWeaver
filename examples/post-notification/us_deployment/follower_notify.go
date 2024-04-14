package main

import (
	"context"
	"fmt"

	"github.com/TiagoMalhadas/xcweaver"
)

// Follower_Notify component.
type Follower_Notify interface {
	Follower_Notify(context.Context, string, int) error
}

// Implementation of the Follower_Notify component.
type follower_Notify struct {
	xcweaver.Implements[Follower_Notify]
	post_storage xcweaver.Ref[Post_storage_america]
}

func (r *follower_Notify) Follower_Notify(ctx context.Context, postId string, userId int) error {

	post, err := r.post_storage.Get().GetPost(ctx, postId)

	if err != nil {
		return err
	}

	fmt.Println("Post: ", post)

	return nil
}
