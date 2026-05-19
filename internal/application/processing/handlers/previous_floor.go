package handlers

import (
	"fmt"
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
)

type PreviousFloorHandler struct {
	config *models.Config
}

func NewPreviousFloorHandler(config *models.Config) *PreviousFloorHandler {
	return &PreviousFloorHandler{config: config}
}

func (h *PreviousFloorHandler) Handle(session *models.PlayerSession, ev events.Event) ([]events.Event, error) {
	e, ok := ev.(events.UserWentToPreviousFloorEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected event type")
	}

	if session.CurrentFloor == 0 {
		return []events.Event{
			events.UserMakesImpossibleMove{
				BaseEvent: events.BaseEvent{Timestamp: e.Timestamp, UserID: e.UserID},
				EventID:   e.GetEventID(),
			},
		}, nil
	}

	session.CurrentFloor--
	session.MonstersKilledOnFloor = 0
	return []events.Event{
		events.UserWentToPreviousFloorEvent{BaseEvent: events.BaseEvent{Timestamp: e.Timestamp, UserID: e.UserID}},
	}, nil
}
