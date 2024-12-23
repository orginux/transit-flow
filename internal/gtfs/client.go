package gtfs

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
	"transit-flow/internal/types"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

// Client handles GTFS realtime data fetching
type Client struct {
	config Config
	client *http.Client
}

// NewClient creates a new GTFS client
func NewClient(config Config) *Client {
	return &Client{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// FetchUpdates fetches and processes GTFS realtime data
func (g *Client) FetchUpdates(ctx context.Context) ([]types.VehicleUpdate, Metrics, error) {
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
func (g *Client) fetchFeed(ctx context.Context) (*gtfs.FeedMessage, error) {
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
