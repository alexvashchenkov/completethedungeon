package parser

import "time"

type RawEvent struct {
	Timestamp time.Time `json:"timestamp"`
	UserID    int       `json:"user_id"`
	EventID   int       `json:"event_id"`
	Extra     string    `json:"extra"`
}
