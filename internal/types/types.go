package types

import "github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"

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
