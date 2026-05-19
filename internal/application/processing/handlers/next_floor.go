package handlers

import (
	"fmt"
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
	"time"
)

type NextFloorHandler struct {
	config *models.Config
}

func NewNextFloorHandler(config *models.Config) *NextFloorHandler {
	return &NextFloorHandler{config: config}
}

func (h *NextFloorHandler) Handle(session *models.PlayerSession, ev events.Event) ([]events.Event, error) {
	e, ok := ev.(events.UserWentToNextFloorEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected event type")
	}

	if h.config == nil || session.MonstersKilledOnFloor < h.config.Monsters {
		return []events.Event{
			events.UserMakesImpossibleMove{
				BaseEvent: events.BaseEvent{Timestamp: e.Timestamp, UserID: e.UserID},
				EventID:   e.GetEventID(),
			},
		}, nil
	}

	if session.CurrentFloor >= h.config.Floors-1 {
		return []events.Event{
			events.UserMakesImpossibleMove{
				BaseEvent: events.BaseEvent{Timestamp: e.Timestamp, UserID: e.UserID},
				EventID:   e.GetEventID(),
			},
		}, nil
	}

	if !session.CurrentFloorStartedAt.IsZero() && !session.CurrentFloorFinishedAt.IsZero() {
		floorDuration := session.CurrentFloorFinishedAt.Sub(session.CurrentFloorStartedAt)
		session.Metrics.FloorDurations = append(session.Metrics.FloorDurations, floorDuration)
	}

	session.CurrentFloor++
	session.MonstersKilledOnFloor = 0
	session.CurrentFloorStartedAt = e.Timestamp
	session.CurrentFloorFinishedAt = time.Time{}
	session.Metrics.FloorsCompleted++

	return []events.Event{
		events.UserWentToNextFloorEvent{BaseEvent: events.BaseEvent{Timestamp: e.Timestamp, UserID: e.UserID}},
	}, nil
}
