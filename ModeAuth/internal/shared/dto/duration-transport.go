package dto

import "time"

type DurationTransport struct {
	UserID    string        `json:"user_id"`
	UserName  string        `json:"user_name"`
	IsBlocked bool          `json:"is_blocked"`
	Duration  time.Duration `json:"duration"`
	UpdatedAt time.Time     `json:"updated_at"`
}
