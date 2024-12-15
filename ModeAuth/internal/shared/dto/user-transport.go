package dto

import "time"

type UserTransport struct {
	ID             string        `json:"id"`
	IsBot          bool          `json:"is_bot"`
	Role           string        `json:"role"`
	UserName       string        `json:"user_name"` // optional
	IsBlocked      bool          `json:"is_blocked"`
	BlockedEarlier int           `json:"blocked_earlier"`
	TimeToUnlock   time.Duration `json:"time_to_unlock"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdateAt       time.Time     `json:"update_at"`
}
