package ycsb

import (
	"time"
)

// MeasurementInfo contains metrics of one measurement.
type MeasurementInfo interface {
	// Get returns the value corresponded to the specified metric, such QPS, MIN, MAX，etc.
	// If metric does not exist, the returned value will be nil.
	Get(metricName string) interface{}
}

// Measurement measures the operations metrics.
type Measurement interface {
	// Measure measures the operation latency.
	Measure(latency time.Duration)
	// Summary returns the summary of the measurement.
	Summary() string
	// Info returns the MeasurementInfo of the measurement.
	Info() MeasurementInfo
}
