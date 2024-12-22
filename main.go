package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
	"google.golang.org/protobuf/proto"
)

// Config holds all configuration parameters
type Config struct {
	FeedURL  string
	Username string
	Password string
	Timeout  time.Duration
}

// GTFSClient handles GTFS realtime data fetching
type GTFSClient struct {
	config Config
	client *http.Client
}

// VehicleUpdate represents processed GTFS data
type VehicleUpdate struct {
	Timestamp int64 `parquet:"name=timestamp, type=INT64, convertedType=TIMESTAMP_MILLIS"`
	Position  *gtfs.Position
	TripID    string `parquet:"name=trip_id,type=BYTE_ARRAY,convertedtype=UTF8"`
	RouteID   string `parquet:"name=route_id,type=BYTE_ARRAY,convertedtype=UTF8"`
	Status    string `parquet:"name=status,type=BYTE_ARRAY,convertedtype=UTF8"`
	StopID    string `parquet:"name=stop_id,type=BYTE_ARRAY,convertedtype=UTF8"`
	Delay     int32  `parquet:"name=delay,type=INT32"`
}

// NewGTFSClient creates a new GTFS client
func NewGTFSClient(config Config) *GTFSClient {
	return &GTFSClient{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// Metrics holds metrics for GTFS data fetching
type Metrics struct {
	FetchTime      time.Duration
	ProcessingTime time.Duration
	TotalTime      time.Duration
	UpdatesCount   int
}

// FetchUpdates fetches and processes GTFS realtime data
func (g *GTFSClient) FetchUpdates(ctx context.Context) ([]VehicleUpdate, Metrics, error) {
	metrics := Metrics{}

	totalStart := time.Now()
	feedStart := time.Now()

	feed, err := g.fetchFeed(ctx)
	if err != nil {
		return nil, metrics, fmt.Errorf("fetch feed: %w", err)
	}
	metrics.FetchTime = time.Since(feedStart)

	if feed == nil {
		return nil, metrics, fmt.Errorf("received nil feed")
	}

	processingStart := time.Now()
	updates := g.processFeed(feed)
	metrics.ProcessingTime = time.Since(processingStart)
	metrics.TotalTime = time.Since(totalStart)
	metrics.UpdatesCount = len(updates)

	return updates, metrics, nil
}

// fetchFeed retrieves raw GTFS data
func (g *GTFSClient) fetchFeed(ctx context.Context) (*gtfs.FeedMessage, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", g.config.FeedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	if g.config.Username != "" && g.config.Password != "" {
		req.SetBasicAuth(g.config.Username, g.config.Password)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	feed := &gtfs.FeedMessage{}
	if err := proto.Unmarshal(body, feed); err != nil {
		return nil, fmt.Errorf("unmarshal proto: %w", err)
	}

	return feed, nil
}

// processTripUpdate handles trip updates
func (g *GTFSClient) processTripUpdate(update *gtfs.TripUpdate, timestamp time.Time) []VehicleUpdate {
	updates := make([]VehicleUpdate, 0, len(update.StopTimeUpdate))

	// Check if Trip is nil
	if update.Trip == nil {
		return updates
	}

	tripID := getString(update.Trip.TripId)
	routeID := getString(update.Trip.RouteId)

	for _, stopTimeUpdate := range update.StopTimeUpdate {
		// Add nil checks for Arrival
		var delay int32
		if stopTimeUpdate.Arrival != nil {
			delay = getInt32(stopTimeUpdate.Arrival.Delay)
		}

		updates = append(updates, VehicleUpdate{
			TripID:    tripID,
			RouteID:   routeID,
			Timestamp: timestamp.UnixMilli(),
			Delay:     delay,
			StopID:    getString(stopTimeUpdate.StopId),
		})
	}

	return updates
}

// processVehicleUpdate handles vehicle updates
func (g *GTFSClient) processVehicleUpdate(vehicle *gtfs.VehiclePosition, timestamp time.Time) *VehicleUpdate {
	if vehicle == nil {
		return nil
	}

	update := &VehicleUpdate{
		Timestamp: timestamp.UnixMilli(),
		Position:  vehicle.Position,
		StopID:    getString(vehicle.StopId),
	}

	if vehicle.Trip != nil {
		update.TripID = getString(vehicle.Trip.TripId)
		update.RouteID = getString(vehicle.Trip.RouteId)
	}

	if vehicle.CurrentStatus != nil {
		update.Status = vehicle.CurrentStatus.String()
	}

	return update
}

// processEntity handles individual GTFS entities
func (g *GTFSClient) processEntity(entity *gtfs.FeedEntity, timestamp time.Time) []VehicleUpdate {
	var updates []VehicleUpdate

	// Check if entity is nil
	if entity == nil {
		return updates
	}

	if entity.TripUpdate != nil {
		updates = append(updates, g.processTripUpdate(entity.TripUpdate, timestamp)...)
	}

	if entity.Vehicle != nil {
		if update := g.processVehicleUpdate(entity.Vehicle, timestamp); update != nil {
			updates = append(updates, *update)
		}
	}

	return updates
}

// processFeed processes GTFS feed data into vehicle updates
func (g *GTFSClient) processFeed(feed *gtfs.FeedMessage) []VehicleUpdate {
	updates := make([]VehicleUpdate, 0)

	// Check if feed or Entity is nil
	if feed == nil || feed.Entity == nil {
		return updates
	}

	timestamp := time.Now()

	for _, entity := range feed.Entity {
		updates = append(updates, g.processEntity(entity, timestamp)...)
	}

	return updates
}

// Helper functions for safer pointer dereferencing
func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func getInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

func writeVehicleUpdates(basePath string, updates []VehicleUpdate) error {

	// Generate filename
	filename := fmt.Sprintf("gtfs_%s.parquet",
		time.Now().Format("2006-01-02_15-04-05"))
	fullPath := filepath.Join(basePath, filename)

	var err error
	fw, err := local.NewLocalFileWriter(fullPath)
	if err != nil {
		return fmt.Errorf("create local file writer: %w", err)
	}
	defer fw.Close()

	pw, err := writer.NewParquetWriter(fw, new(VehicleUpdate), 4)
	if err != nil {
		return fmt.Errorf("create parquet writer: %w", err)
	}

	// pw.RowGroupSize = 128 * 1024 * 1024 //128M
	// pw.PageSize = 8 * 1024              //8K
	pw.CompressionType = parquet.CompressionCodec_ZSTD

	for _, vu := range updates {
		if err = pw.Write(vu); err != nil {
			log.Println("Write error", err)
		}
	}

	if err = pw.WriteStop(); err != nil {
		return fmt.Errorf("write stop: %w", err)
	}
	fw.Close()
	return nil
}

func main() {
	config := Config{
		FeedURL: "https://zet.hr/gtfs-rt-protobuf",
		Timeout: 10 * time.Second,
	}

	client := NewGTFSClient(config)
	ctx := context.Background()

	updates, metrics, err := client.FetchUpdates(ctx)
	if err != nil {
		log.Fatalf("Failed to fetch updates: %v", err)
	}

	fmt.Printf("Fetched %d updates in %v\n", metrics.UpdatesCount, metrics.TotalTime)

	if err = writeVehicleUpdates("output/", updates); err != nil {
		log.Fatalf("Failed to write updates: %v", err)
	}
}
