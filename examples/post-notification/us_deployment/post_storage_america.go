package main

import (
	"context"

	"github.com/TiagoMalhadas/xcweaver"
)

// Post_storage_america component.
type Post_storage_america interface {
	GetPost(context.Context, string) (string, error)
}

// Implementation of the Post_storage_america component.
type post_storage_america struct {
	xcweaver.Implements[Post_storage_america]
	clientRedis xcweaver.Antipode
}

func (r *post_storage_america) GetPost(ctx context.Context, postId string) (string, error) {

	err := r.clientRedis.Barrier(ctx)
	if err != nil {
		return "error on barrier", err
	}

	post, _, err := r.clientRedis.Read(ctx, "posts", postId)
	if err != nil {
		inconsistencies.Inc()
		return "post not found", nil
	}

	return post, nil
}
