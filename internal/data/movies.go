package data

import (
	"time"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // exclude from JSON output
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"` // Runtime type found in internal/data/runtime.go
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}
