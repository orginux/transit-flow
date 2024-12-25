package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"transit-flow/internal/gtfs"
	"transit-flow/internal/storage"
	"transit-flow/internal/types"
)

// Storage defines interface for both local and GCS storage
type Storage interface {
	Write(ctx context.Context, updates []types.VehicleUpdate) (string, error)
}

// NewStorage creates a new storage provider based on the environment
func NewStorage(config storage.Config) Storage {
	if isLocalOnly() {
		return storage.NewLocalStorage(config)
	}
	return storage.NewGCSStorage(config)
}

// isLocalOnly returns true if GTFS_LOCAL_ONLY environment variable is set
func isLocalOnly() bool {
	_, exists := os.LookupEnv("GTFS_LOCAL_ONLY")
	return exists
}

func main() {
	// Start minimal HTTP server for health check
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "OK")
		})
		http.ListenAndServe(":"+port, nil)
	}()

	// Fetch GTFS updates
	config := gtfs.Config{
		FeedURL: "https://zet.hr/gtfs-rt-protobuf",
		Timeout: 10 * time.Second,
	}
	client := gtfs.NewClient(config)
	ctx := context.Background()
	updates, metrics, err := client.FetchUpdates(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch updates: %v", err)
	}
	fmt.Printf("Fetched %d updates in %v\n", metrics.UpdatesCount, metrics.TotalTime)

	// Write updates to storage
	storageConfig := storage.Config{
		BasePath:   "output/",
		TimeFormat: "2006-01-02_15-04-05.000",
		FilePrefix: "gtfs",
		BucketName: os.Getenv("GTFS_GSC_BUCKET"),
		ProjectID:  os.Getenv("GTFS_GSC_PROJECT"),
	}
	storageProvider := NewStorage(storageConfig)
	path, err := storageProvider.Write(ctx, updates)
	if err != nil {
		log.Fatalf("Failed to write updates: %v", err)
	}
	fmt.Printf("Wrote updates to %s\n", path)
}
