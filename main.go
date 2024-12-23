package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"transit-flow/internal/gtfs"
	"transit-flow/internal/storage"
)

func main() {
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

	localConfig := storage.Config{
		BasePath:   "output/",
		TimeFormat: "2006-01-02_15-04-05.000",
		FilePrefix: "gtfs",
	}
	localStorage := storage.NewLocalStorage(localConfig)

	storageProvider := localStorage
	path, err := storageProvider.Write(ctx, updates)
	if err != nil {
		log.Fatalf("Failed to write updates: %v", err)
	}
	fmt.Printf("Wrote updates to %s\n", path)
}
