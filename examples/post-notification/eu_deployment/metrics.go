package main

import "github.com/TiagoMalhadas/xcweaver/metrics"

type RegionLabel struct {
	Region string
}

var (
	notificationsSent = metrics.NewCounter(
		"sn_notificationsSent",
		"The number of notifications sent over the queue",
	)
)
