package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/TiagoMalhadas/xcweaver"
)

func main() {
	if err := xcweaver.Run(context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

type app struct {
	xcweaver.Implements[xcweaver.Main]
	post_upload       xcweaver.Ref[Post_upload]
	post_notification xcweaver.Listener
}

// serve is called by xcweaver.Run and contains the body of the application.
func serve(ctx context.Context, app *app) error {

	fmt.Printf("post-notification listener available on %v\n", app.post_notification)

	var post_upload Post_upload = app.post_upload.Get()

	// Serve the /post_notification endpoint.
	http.HandleFunc("/post_notification", func(w http.ResponseWriter, r *http.Request) {
		post := r.URL.Query().Get("post")
		err := post_upload.Post(ctx, post, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	return http.Serve(app.post_notification, nil)
}
