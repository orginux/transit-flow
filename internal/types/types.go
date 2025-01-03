package types

// VehicleUpdate represents raw GTFS Realtime data
type VehicleUpdate struct {
	// Common timestamp for event
	Timestamp int64 `parquet:"name=timestamp, type=INT64, convertedType=TIMESTAMP_MILLIS"`

	// Trip information
	TripID      string `parquet:"name=trip_id,type=BYTE_ARRAY,convertedtype=UTF8"`
	RouteID     string `parquet:"name=route_id,type=BYTE_ARRAY,convertedtype=UTF8"`
	DirectionID int32  `parquet:"name=direction_id,type=INT32,convertedtype=UINT_32"`

	// Vehicle information
	VehicleID string `parquet:"name=vehicle_id,type=BYTE_ARRAY,convertedtype=UTF8"`

	// Position data
	Latitude  float32 `parquet:"name=latitude, type=FLOAT"`
	Longitude float32 `parquet:"name=longitude, type=FLOAT"`
	Bearing   float32 `parquet:"name=bearing, type=FLOAT"`
	Speed     float32 `parquet:"name=speed, type=FLOAT"`

	// Stop data
	StopID       string `parquet:"name=stop_id,type=BYTE_ARRAY,convertedtype=UTF8"`
	StopSequence int32  `parquet:"name=stop_sequence,type=INT32,convertedtype=UINT_32"`
	Status       string `parquet:"name=status,type=BYTE_ARRAY,convertedtype=UTF8"`
	Delay        int32  `parquet:"name=delay,type=INT32"`
	ArrivalTime  int64  `parquet:"name=arrival_time,type=INT64,convertedtype=TIMESTAMP_MILLIS"`
}
