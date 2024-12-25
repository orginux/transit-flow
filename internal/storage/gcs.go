package storage

import (
	"context"
	"fmt"
	"transit-flow/internal/types"

	"github.com/xitongsys/parquet-go-source/gcs"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

type GCSStorage struct {
	config Config
}

func NewGCSStorage(config Config) *GCSStorage {
	return &GCSStorage{
		config: config,
	}
}

func (gs *GCSStorage) Write(ctx context.Context, updates []types.VehicleUpdate) (string, error) {
	filename := filenameGenerator(gs.config)
	fullPath := fmt.Sprintf("gs://%s/%s/%s", gs.config.BucketName, gs.config.BasePath, filename)

	// Create a new GCS file writer
	fw, err := gcs.NewGcsFileWriter(
		ctx,
		gs.config.ProjectID,
		gs.config.BucketName,
		gs.config.BasePath+filename,
	)
	if err != nil {
		return "", fmt.Errorf("create GCS file writer: %w", err)
	}
	defer fw.Close()

	// Create a new parquet writer with 4 go routines
	pw, err := writer.NewParquetWriter(fw, new(types.VehicleUpdate), 4)
	if err != nil {
		return "", fmt.Errorf("create parquet writer: %w", err)
	}

	// Set compression type to ZSTD (same as in local storage)
	pw.CompressionType = parquet.CompressionCodec_ZSTD

	// Write all updates to the parquet file
	for _, vu := range updates {
		if err = pw.Write(vu); err != nil {
			return "", fmt.Errorf("write record: %w", err)
		}
	}

	// Finalize writing
	if err = pw.WriteStop(); err != nil {
		return "", fmt.Errorf("write stop: %w", err)
	}

	return fullPath, nil
}
