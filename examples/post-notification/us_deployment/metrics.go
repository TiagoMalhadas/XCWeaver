package main

import "github.com/TiagoMalhadas/xcweaver/metrics"

var (
	inconsistencies = metrics.NewCounter(
		"sn_inconsistencies",
		"The number of times an cross-service inconsistency has occured",
	)
	notificationsReceived = metrics.NewCounter(
		"notificationsReceived",
		"The number of notifications received",
	)
	postNotificationDuration = metrics.NewHistogram(
		"sn_post_notification_duration_ms",
		"Duration of post-notification requests in milliseconds",
		metrics.NonNegativeBuckets,
	)
	readPostDurationMs = metrics.NewHistogram(
		"sn_read_post_duration_ms",
		"Duration of read operation in milliseconds in the us region",
		metrics.NonNegativeBuckets,
	)
	queueDurationMs = metrics.NewHistogram(
		"sn_queue_duration_ms",
		"Duration of queue in milliseconds",
		metrics.NonNegativeBuckets,
	)
)
