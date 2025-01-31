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

	// Get trip-level information
	tripID := getString(update.Trip.TripId)
	routeID := getString(update.Trip.RouteId)
	direction := int32(update.GetTrip().GetDirectionId())
	startTime := getString(update.Trip.StartTime)
	startDate := getString(update.Trip.StartDate)
	scheduleRelationship := update.Trip.GetScheduleRelationship().String()

	// Get vehicle information
	var vehicleID, vehicleLabel string
	if update.Vehicle != nil {
		vehicleID = getString(update.Vehicle.Id)
		vehicleLabel = getString(update.Vehicle.Label)
	}

	for _, stopTimeUpdate := range update.StopTimeUpdate {
		var arrivalDelay int32
		var departureDelay int32
		var arrivalTime time.Time
		var departureTime time.Time
		var arrivalUncertainty int32
		var departureUncertainty int32

		// Process arrival information
		if stopTimeUpdate.Arrival != nil {
			arrivalDelay = getInt32(stopTimeUpdate.Arrival.Delay)
			arrivalUncertainty = getInt32(stopTimeUpdate.Arrival.Uncertainty)
			if stopTimeUpdate.Arrival.Time != nil {
				arrivalTime = time.Unix(*stopTimeUpdate.Arrival.Time, 0)
			}
		}

		// Process departure information
		if stopTimeUpdate.Departure != nil {
			departureDelay = getInt32(stopTimeUpdate.Departure.Delay)
			departureUncertainty = getInt32(stopTimeUpdate.Departure.Uncertainty)
			if stopTimeUpdate.Departure.Time != nil {
				departureTime = time.Unix(*stopTimeUpdate.Departure.Time, 0)
			}
		}

		updates = append(updates, types.VehicleUpdate{
			// Trip information
			TripID:               tripID,
			RouteID:              routeID,
			DirectionID:          direction,
			StartTime:            startTime,
			StartDate:            startDate,
			ScheduleRelationship: scheduleRelationship,

			// Vehicle information
			VehicleID:    vehicleID,
			VehicleLabel: vehicleLabel,

			// Stop information
			StopID:       getString(stopTimeUpdate.StopId),
			StopSequence: int32(getUint32(stopTimeUpdate.StopSequence)),

			// Arrival information
			ArrivalTime:        arrivalTime.UnixMilli(),
			ArrivalDelay:       arrivalDelay,
			ArrivalUncertainty: arrivalUncertainty,

			// Departure information
			DepartureTime:        departureTime.UnixMilli(),
			DepartureDelay:       departureDelay,
			DepartureUncertainty: departureUncertainty,

			// Common timestamp
			Timestamp: timestamp.UnixMilli(),
		})
	}
	return updates
}

// processVehicleUpdate handles vehicle updates
func (g *Client) processVehicleUpdate(vehicle *gtfs.VehiclePosition, timestamp time.Time) *types.VehicleUpdate {
	if vehicle == nil || vehicle.Position == nil {
		return nil
	}

	update := &types.VehicleUpdate{
		Timestamp: timestamp.UnixMilli(),

		// Position data
		Latitude:  vehicle.Position.GetLatitude(),
		Longitude: vehicle.Position.GetLongitude(),
		Bearing:   vehicle.Position.GetBearing(),
		Speed:     vehicle.Position.GetSpeed(),
		StopID:    getString(vehicle.StopId),
	}

	// Process trip information
	if vehicle.Trip != nil {
		update.TripID = getString(vehicle.Trip.TripId)
		update.RouteID = getString(vehicle.Trip.RouteId)
		update.DirectionID = int32(vehicle.Trip.GetDirectionId())
		update.StartTime = getString(vehicle.Trip.StartTime)
		update.StartDate = getString(vehicle.Trip.StartDate)
		update.ScheduleRelationship = vehicle.Trip.GetScheduleRelationship().String()
	}

	// Process stop information
	if vehicle.CurrentStatus != nil {
		update.Status = vehicle.CurrentStatus.String()
	}
	if vehicle.CurrentStopSequence != nil {
		update.StopSequence = int32(vehicle.GetCurrentStopSequence())
	}

	// Process vehicle information
	if vehicle.Vehicle != nil {
		update.VehicleID = getString(vehicle.Vehicle.Id)
		update.VehicleLabel = getString(vehicle.Vehicle.Label)
	}

	// Process congestion and occupancy information
	if vehicle.CongestionLevel != nil {
		update.CongestionLevel = vehicle.CongestionLevel.String()
	}
	if vehicle.OccupancyStatus != nil {
		update.OccupancyStatus = vehicle.OccupancyStatus.String()
	}
	if vehicle.OccupancyPercentage != nil {
		update.OccupancyPercentage = int32(vehicle.GetOccupancyPercentage())
	}

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

// ProcessFeed processes GTFS feed data into updates
func (g *Client) ProcessFeed(feed *gtfs.FeedMessage) []types.VehicleUpdate {
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

func getUint32(i *uint32) uint32 {
	if i == nil {
		return 0
	}
	return *i
}
