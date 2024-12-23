package storage

import (
	"context"
	"fmt"
	"path/filepath"
	"transit-flow/internal/types"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

type LocalStorage struct {
	config Config
}

func NewLocalStorage(config Config) *LocalStorage {
	return &LocalStorage{
		config: config,
	}
}

func (ls *LocalStorage) Write(ctx context.Context, updates []types.VehicleUpdate) (string, error) {
	filename := filenameGenerator(ls.config)
	fullPath := filepath.Join(ls.config.BasePath, filename)

	fw, err := local.NewLocalFileWriter(fullPath)
	if err != nil {
		return "", fmt.Errorf("create local file writer: %w", err)
	}
	defer fw.Close()

	pw, err := writer.NewParquetWriter(fw, new(types.VehicleUpdate), 4)
	if err != nil {
		return "", fmt.Errorf("create parquet writer: %w", err)
	}

	pw.CompressionType = parquet.CompressionCodec_ZSTD

	for _, vu := range updates {
		if err = pw.Write(vu); err != nil {
			return "", fmt.Errorf("write record: %w", err)
		}
	}

	if err = pw.WriteStop(); err != nil {
		return "", fmt.Errorf("write stop: %w", err)
	}

	return fullPath, nil
}
