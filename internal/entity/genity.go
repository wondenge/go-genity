package entity

import (
	"time"
)

// Genity represents an genity record.
type Genity struct {
	Title     string    `json:"Title"`
	ID        string    `json:"id"`
	Timestamp time.Time `json:"Timestamp"`
}
