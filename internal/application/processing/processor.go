package processing

import (
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
)

type Processor struct {
}

func (p *Processor) Process(session *models.PlayerSession, event events.Event) ([]events.Event, error) {
	return []events.Event{event}, nil
}
