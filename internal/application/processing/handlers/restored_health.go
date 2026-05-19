package handlers

import (
	"fmt"
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
)

type RestoredHealthHandler struct{}

func NewRestoredHealthHandler() *RestoredHealthHandler {
	return &RestoredHealthHandler{}
}

func (h *RestoredHealthHandler) Handle(session *models.PlayerSession, ev events.Event) ([]events.Event, error) {
	e, ok := ev.(events.UserRestoredHealthEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected event type")
	}

	session.HP += e.Amount
	if session.HP > 100 {
		session.HP = 100
	}

	return []events.Event{
		events.UserRestoredHealthEvent{BaseEvent: events.BaseEvent{Timestamp: e.Timestamp, UserID: e.UserID}, Amount: e.Amount},
	}, nil
}
