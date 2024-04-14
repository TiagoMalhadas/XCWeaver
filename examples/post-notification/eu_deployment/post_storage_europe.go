package main

import (
	"context"
	"fmt"

	"github.com/TiagoMalhadas/xcweaver"
	"github.com/google/uuid"
)

// Post_storage component.
type Post_storage_europe interface {
	Post(context.Context, string) ([]byte, string, error)
}

// Implementation of the Post_storage component.
type post_storage_europe struct {
	xcweaver.Implements[Post_storage_europe]
	clientRedis xcweaver.Antipode
}

func (r *post_storage_europe) Post(ctx context.Context, post string) ([]byte, string, error) {

	id := uuid.New()

	ctx, err := r.clientRedis.Write(ctx, "posts", id.String(), post)
	
	
	fmt.Println("plaaa")
	if err != nil {
		return []byte{}, "", err
	}

	lineage, err := xcweaver.GetLineage(ctx)

	if err != nil {
		return []byte{}, "", err
	}

	return lineage, id.String(), nil
}
