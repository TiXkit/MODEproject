package models

import "time"

type UserTransport struct {
	ID             string
	IsBot          bool
	Role           string
	UserName       string
	IsBlocked      bool
	BlockedEarlier int
	TimeToUnlock   time.Duration
	CreatedAt      time.Time
	UpdateAt       time.Time
}
