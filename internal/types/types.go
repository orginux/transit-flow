package types

// VehicleUpdate represents raw GTFS Realtime data
type VehicleUpdate struct {
	// Common timestamp for event
	Timestamp int64 `parquet:"name=timestamp,type=INT64,convertedType=TIMESTAMP_MILLIS"`

	// Trip information
	TripID               string `parquet:"name=trip_id,type=BYTE_ARRAY,convertedtype=UTF8"`
	RouteID              string `parquet:"name=route_id,type=BYTE_ARRAY,convertedtype=UTF8"`
	DirectionID          int32  `parquet:"name=direction_id,type=INT32"`
	StartTime            string `parquet:"name=start_time,type=BYTE_ARRAY,convertedtype=UTF8"`
	StartDate            string `parquet:"name=start_date,type=BYTE_ARRAY,convertedtype=UTF8"`
	ScheduleRelationship string `parquet:"name=schedule_relationship,type=BYTE_ARRAY,convertedtype=UTF8"`

	// Vehicle information
	VehicleID    string `parquet:"name=vehicle_id,type=BYTE_ARRAY,convertedtype=UTF8"`
	VehicleLabel string `parquet:"name=vehicle_label,type=BYTE_ARRAY,convertedtype=UTF8"`

	// Position data
	Latitude  float32 `parquet:"name=latitude,type=FLOAT"`
	Longitude float32 `parquet:"name=longitude,type=FLOAT"`
	Bearing   float32 `parquet:"name=bearing,type=FLOAT"`
	Speed     float32 `parquet:"name=speed,type=FLOAT"`

	// Stop data
	StopID       string `parquet:"name=stop_id,type=BYTE_ARRAY,convertedtype=UTF8"`
	StopSequence int32  `parquet:"name=stop_sequence,type=INT32"`                  // Changed back to int32
	Status       string `parquet:"name=status,type=BYTE_ARRAY,convertedtype=UTF8"` // For VehicleStopStatus (INCOMING_AT, STOPPED_AT, IN_TRANSIT_TO)

	// Arrival information
	ArrivalTime        int64 `parquet:"name=arrival_time,type=INT64,convertedType=TIMESTAMP_MILLIS"`
	ArrivalDelay       int32 `parquet:"name=arrival_delay,type=INT32"`
	ArrivalUncertainty int32 `parquet:"name=arrival_uncertainty,type=INT32"`

	// Departure information
	DepartureTime        int64 `parquet:"name=departure_time,type=INT64,convertedType=TIMESTAMP_MILLIS"`
	DepartureDelay       int32 `parquet:"name=departure_delay,type=INT32"`
	DepartureUncertainty int32 `parquet:"name=departure_uncertainty,type=INT32"`

	// Vehicle congestion and occupancy (optional)
	CongestionLevel     string `parquet:"name=congestion_level,type=BYTE_ARRAY,convertedtype=UTF8"` // From VehiclePosition_CongestionLevel
	OccupancyStatus     string `parquet:"name=occupancy_status,type=BYTE_ARRAY,convertedtype=UTF8"` // From VehiclePosition_OccupancyStatus
	OccupancyPercentage int32  `parquet:"name=occupancy_percentage,type=INT32"`
}
