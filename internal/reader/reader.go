// Package reader provides utilities for reading and parsing
// newline-delimited JSON (NDJSON) log streams from various sources.
package reader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

// Record represents a single parsed JSON log entry.
type Record map[string]interface{}

// Reader reads newline-delimited JSON records from an io.Reader.
type Reader struct {
	scanner *bufio.Scanner
	source  string
}

// New creates a new Reader wrapping the given io.Reader.
// The source label is attached to each record under the "_source" key.
func New(r io.Reader, source string) *Reader {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	return &Reader{
		scanner: scanner,
		source:  source,
	}
}

// Next reads the next JSON record from the stream.
// It returns (nil, io.EOF) when the stream is exhausted.
// Non-JSON lines are skipped with a parse error returned.
func (r *Reader) Next() (Record, error) {
	for r.scanner.Scan() {
		line := r.scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var rec Record
		if err := json.Unmarshal(line, &rec); err != nil {
			return nil, fmt.Errorf("parse error in %q: %w", r.source, err)
		}
		if r.source != "" {
			rec["_source"] = r.source
		}
		return rec, nil
	}
	if err := r.scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan error in %q: %w", r.source, err)
	}
	return nil, io.EOF
}

// ReadAll reads all records from the stream, skipping lines that fail to parse.
// Errors encountered per line are collected and returned alongside the records.
func (r *Reader) ReadAll() ([]Record, []error) {
	var records []Record
	var errs []error
	for {
		rec, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			errs = append(errs, err)
			continue
		}
		records = append(records, rec)
	}
	return records, errs
}
