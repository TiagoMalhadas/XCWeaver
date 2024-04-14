package main

import (
	"context"
	"log"

	"github.com/TiagoMalhadas/xcweaver"
)

func main() {
	if err := xcweaver.Run(context.Background(), serve); err != nil {
		log.Fatal(err)
	}
}

type app struct {
	xcweaver.Implements[xcweaver.Main]
	notifier xcweaver.Ref[Notifier]
}

// serve is called by xcweaver.Run and contains the body of the application.
func serve(ctx context.Context, app *app) error {

	i := 0
	for i == 0 {

	}

	return nil
}
