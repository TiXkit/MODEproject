package dto

import "time"

type Report struct {
	ID            int       `json:"id"`
	UserID        string    `json:"user_id"`
	Status        bool      `json:"status"`
	TimeSending   time.Time `json:"time_sending"`
	ChatPhoto     string    `json:"chat_photo"`
	FeedBackPhoto string    `json:"feed_back_photo"`
	Comment       string    `json:"comment"`
}
