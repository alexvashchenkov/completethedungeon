package handlers

import (
	"fmt"
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
	"time"
)

type EnteredBossFloorHandler struct {
	config *models.Config
}

func NewEnteredBossFloorHandler(config *models.Config) *EnteredBossFloorHandler {
	return &EnteredBossFloorHandler{config: config}
}
func (h *EnteredBossFloorHandler) Handle(session *models.PlayerSession, ev events.Event) ([]events.Event, error) {
	e, ok := ev.(events.UserEnteredBossFloorEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected event type")
	}

	if session.CurrentFloor >= h.config.Floors {
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

	session.CurrentFloor = h.config.Floors
	session.CurrentFloorStartedAt = e.Timestamp
	session.CurrentFloorFinishedAt = time.Time{}
	session.MonstersKilledOnFloor = 0

	return []events.Event{
		events.UserEnteredBossFloorEvent{BaseEvent: events.BaseEvent{Timestamp: e.Timestamp, UserID: e.UserID}},
	}, nil
}
