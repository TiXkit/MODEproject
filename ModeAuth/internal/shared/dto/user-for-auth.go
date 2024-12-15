package dto

import "time"

type User struct {
	ID              string    `json:"id"`
	Role            string    `json:"role"`
	UserName        string    `json:"user_name"` // optional
	ApprovedReports int       `json:"approved_reports"`
	OrderTaken      int       `json:"order_taken"`
	BlockedEarlier  int       `json:"blocked_earlier"`
	CreatedAt       time.Time `json:"created_at"`
	UpdateAt        time.Time `json:"update_at"`
}
