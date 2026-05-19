package output

import "impulse/internal/domain/models/events"

type EventSink interface {
	Write(events.Event) error
	WriteMany([]events.Event) error
}
