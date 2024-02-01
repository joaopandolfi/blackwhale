package compressible

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
)

type Compress interface {
	Zip() error
	Unzip() error
}

type Compressible struct {
	Compressed bool `json:"compressed"`
}

func ZipBytes(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("creating compressor: %w", err)
	}
	if _, err := gz.Write(data); err != nil {
		return nil, fmt.Errorf("compressing data: %w", err)
	}
	if err := gz.Flush(); err != nil {
		return nil, fmt.Errorf("flushing gzip buffer: %w", err)
	}
	if err := gz.Close(); err != nil {
		return nil, fmt.Errorf("closing virtual file: %w", err)
	}

	return b.Bytes(), nil
}

func UnzipBytes(data []byte) ([]byte, error) {
	rdata := bytes.NewReader(data)

	r, err := gzip.NewReader(rdata)
	if err != nil {
		return nil, fmt.Errorf("decompressing data: %w", err)
	}

	s, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading buffer: %w", err)
	}

	return s, nil
}

func (m *Compressible) Compress(vals []*string) error {
	if m.Compressed {
		return nil
	}

	for i, val := range vals {
		compressedVal, err := ZipBytes([]byte(*val))
		if err != nil {
			return fmt.Errorf("compressing %d: %w", i, err)
		}
		*vals[i] = base64.StdEncoding.EncodeToString(compressedVal)
	}
	m.Compressed = true
	return nil
}

func (m *Compressible) Decompress(vals []*string) error {
	if !m.Compressed {
		return nil
	}

	for i, val := range vals {
		data, err := base64.StdEncoding.DecodeString(*val)
		if err != nil {
			return fmt.Errorf("decoding compressed value: %d: %w", i, err)
		}
		uncompressedVal, err := UnzipBytes(data)
		if err != nil {
			return fmt.Errorf("uncompressing data: %d : %w", i, err)
		}
		*vals[i] = string(uncompressedVal)
	}
	m.Compressed = false
	return nil
}
