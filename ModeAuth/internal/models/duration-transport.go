package models

import "time"

type DurationTransport struct {
	UserID    string
	UserName  string
	IsBlocked bool
	Duration  time.Duration
	UpdatedAt time.Time
}
