package types

// VehicleUpdate represents processed GTFS data
type VehicleUpdate struct {
	Timestamp    int64   `parquet:"name=timestamp, type=INT64, convertedType=TIMESTAMP_MILLIS"`
	VehicleID    string  `parquet:"name=vehicle_id,type=BYTE_ARRAY,convertedtype=UTF8"`
	Latitude     float32 `parquet:"name=latitude, type=FLOAT"`
	Longitude    float32 `parquet:"name=longitude, type=FLOAT"`
	Bearing      float32 `parquet:"name=bearing, type=FLOAT"`
	TripID       string  `parquet:"name=trip_id,type=BYTE_ARRAY,convertedtype=UTF8"`
	RouteID      string  `parquet:"name=route_id,type=BYTE_ARRAY,convertedtype=UTF8"`
	Status       string  `parquet:"name=status,type=BYTE_ARRAY,convertedtype=UTF8"`
	StopID       string  `parquet:"name=stop_id,type=BYTE_ARRAY,convertedtype=UTF8"`
	StopSequence int32   `parquet:"name=stop_sequence,type=INT32,convertedtype=UINT_32"`
	Delay        int32   `parquet:"name=delay,type=INT32"`
}
