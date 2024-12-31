package gtfs

import (
	"time"
	"transit-flow/internal/types"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
)

// processTripUpdate handles trip updates
func (g *Client) processTripUpdate(update *gtfs.TripUpdate, timestamp time.Time) []types.VehicleUpdate {
	updates := make([]types.VehicleUpdate, 0, len(update.StopTimeUpdate))

	if update.Trip == nil {
		return updates
	}

	tripID := getString(update.Trip.TripId)
	routeID := getString(update.Trip.RouteId)

	for _, stopTimeUpdate := range update.StopTimeUpdate {
		var delay int32
		if stopTimeUpdate.Arrival != nil {
			delay = getInt32(stopTimeUpdate.Arrival.Delay)
		}

		updates = append(updates, types.VehicleUpdate{
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
func (g *Client) processVehicleUpdate(vehicle *gtfs.VehiclePosition, timestamp time.Time) *types.VehicleUpdate {
	if vehicle == nil {
		return nil
	}

	update := &types.VehicleUpdate{
		Timestamp: timestamp.UnixMilli(),
		StopID:    getString(vehicle.StopId),
	}

	if vehicle.Trip != nil {
		update.TripID = getString(vehicle.Trip.TripId)
		update.RouteID = getString(vehicle.Trip.RouteId)
	}

	if vehicle.CurrentStatus != nil {
		update.Status = vehicle.CurrentStatus.String()
	}

	if vehicle.CurrentStopSequence != nil {
		seq := vehicle.GetCurrentStopSequence()
		update.StopSequence = int32(seq)
	}

	if vehicle.Vehicle.Id != nil {
		update.VehicleID = getString(vehicle.GetVehicle().Id)
	}

	update.Latitude = vehicle.Position.GetLatitude()
	update.Longitude = vehicle.Position.GetLongitude()
	update.Bearing = vehicle.Position.GetBearing()

	return update
}

// processEntity handles individual GTFS entities
func (g *Client) processEntity(entity *gtfs.FeedEntity, timestamp time.Time) []types.VehicleUpdate {
	var updates []types.VehicleUpdate

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
func (g *Client) processFeed(feed *gtfs.FeedMessage) []types.VehicleUpdate {
	updates := make([]types.VehicleUpdate, 0)

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
