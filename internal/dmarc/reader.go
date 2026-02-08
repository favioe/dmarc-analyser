package dmarc

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// OpenReport opens a file (path) and returns a reader that yields decompressed XML.
// Supports .gz (gzip) and .zip (reads first .xml entry). Caller must close the returned reader.
func OpenReport(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".gz":
		gz, err := gzip.NewReader(f)
		if err != nil {
			f.Close()
			return nil, fmt.Errorf("gzip: %w", err)
		}
		return &gzipReadCloser{file: f, gz: gz}, nil
	case ".zip":
		info, err := f.Stat()
		if err != nil {
			f.Close()
			return nil, err
		}
		zr, err := zip.NewReader(f, info.Size())
		if err != nil {
			f.Close()
			return nil, fmt.Errorf("zip: %w", err)
		}
		for _, e := range zr.File {
			if strings.HasSuffix(strings.ToLower(e.Name), ".xml") {
				rc, err := e.Open()
				if err != nil {
					f.Close()
					return nil, err
				}
				return &zipEntryReadCloser{file: f, entry: rc}, nil
			}
		}
		f.Close()
		return nil, fmt.Errorf("zip: no .xml entry found")
	default:
		return f, nil
	}
}

type gzipReadCloser struct {
	file *os.File
	gz   *gzip.Reader
}

func (g *gzipReadCloser) Read(p []byte) (n int, err error) {
	return g.gz.Read(p)
}

func (g *gzipReadCloser) Close() error {
	_ = g.gz.Close()
	return g.file.Close()
}

type zipEntryReadCloser struct {
	file  *os.File
	entry io.ReadCloser
}

func (z *zipEntryReadCloser) Read(p []byte) (n int, err error) {
	return z.entry.Read(p)
}

func (z *zipEntryReadCloser) Close() error {
	_ = z.entry.Close()
	return z.file.Close()
}
