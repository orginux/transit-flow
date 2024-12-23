package gtfs

import "time"

// Config holds all configuration parameters
type Config struct {
	FeedURL  string
	Username string
	Password string
	Timeout  time.Duration
}

// Metrics holds metrics for GTFS data fetching
type Metrics struct {
	FetchTime      time.Duration
	ProcessingTime time.Duration
	TotalTime      time.Duration
	UpdatesCount   int
}
