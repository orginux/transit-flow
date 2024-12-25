package storage

import (
	"context"
	"fmt"
	"time"
	"transit-flow/internal/types"
)

// StorageProvider defines the interface for storage operations
type StorageProvider interface {
	// Write writes data to storage and returns the full path/URL
	Write(ctx context.Context, data []types.VehicleUpdate) (string, error)
}

// Config holds common storage configuration
type Config struct {
	BasePath   string
	TimeFormat string
	FilePrefix string
	BucketName string // GCS bucket name
	ProjectID  string // GCS project ID
}

func filenameGenerator(config Config) string {
	filename := fmt.Sprintf("gtfs_%s.parquet",
		time.Now().Format(config.TimeFormat))
	return filename
}
