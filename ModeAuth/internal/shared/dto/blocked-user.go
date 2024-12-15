package dto

import "time"

type BlockedUser struct {
	ID        string    `json:"id"`
	BlockedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"update_at"`
}
