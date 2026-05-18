package processing

import (
	"errors"
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
)

type Processor struct {
	handlers map[events.EventID]EventHandler
}

func NewProcessor() *Processor {
	return &Processor{
		handlers: make(map[events.EventID]EventHandler),
	}
}

func (p *Processor) Process(session *models.PlayerSession, event events.Event) ([]events.Event, error) {
	handler, ok := p.handlers[event.ID()]
	if !ok {
		return nil, errors.New("unknown event id")
	}

	return handler.Handle(session, event)
}
