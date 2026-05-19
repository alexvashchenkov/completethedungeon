package events

import "impulse/internal/domain/models/events"

type EventSource interface {
	Next() (events.Event, error)
}
