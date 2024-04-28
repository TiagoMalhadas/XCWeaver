package main

import "github.com/TiagoMalhadas/xcweaver/metrics"

var (
	notificationsSent = metrics.NewCounter(
		"sn_notificationsSent",
		"The number of notifications sent over the queue",
	)
	requests = metrics.NewCounter(
		"sn_requests",
		"The number of post-notification requests",
	)
	writePostDurationMs = metrics.NewHistogram(
		"sn_write_post_duration_ms",
		"Duration of write operation in milliseconds in the eu region",
		metrics.NonNegativeBuckets,
	)
)
