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

func (p *post_storage_america) Init(ctx context.Context) error {
	logger := p.Logger(ctx)

	logger.Info("compose post service at us running!")
	return nil
}

func (p *post_storage_america) GetPost(ctx context.Context, postId string) (string, error) {
	logger := p.Logger(ctx)

	logger.Debug("calling barrier!")
	err := p.clientRedis.Barrier(ctx)
	if err != nil {
		logger.Error("error on barrier", "msg", err.Error())
		return "error on barrier", err
	}
	logger.Debug("barrier executed successfully!")

	post, _, err := p.clientRedis.Read(ctx, "posts", postId)
	if err != nil {
		inconsistencies.Inc()
		logger.Error("post not found")
		return "post not found", nil
	}

	return post, nil
}
