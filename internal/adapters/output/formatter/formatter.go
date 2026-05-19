package formatter

import "impulse/internal/domain/models/events"

type Formatter interface {
	Format(events.Event) string
}
