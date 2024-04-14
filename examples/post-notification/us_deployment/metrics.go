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
)
