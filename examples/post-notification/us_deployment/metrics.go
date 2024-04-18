package main

import "github.com/TiagoMalhadas/xcweaver/metrics"

type RegionLabel struct {
	Region string
}

var (
	inconsistencies = metrics.NewCounter(
		"sn_inconsistencies",
		"The number of times an cross-service inconsistency has occured",
	)
	notificationsReceived = metrics.NewCounter(
		"notificationsReceived",
		"The number of notifications received",
	)
)
