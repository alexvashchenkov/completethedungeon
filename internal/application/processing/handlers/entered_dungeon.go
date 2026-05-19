package handlers

import (
	"fmt"
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
	"time"
)

type EnteredDungeonHandler struct {
	config *models.Config
}

func NewEnteredDungeonHandler(config *models.Config) *EnteredDungeonHandler {
	return &EnteredDungeonHandler{config: config}
}

func (h *EnteredDungeonHandler) Handle(session *models.PlayerSession, ev events.Event) ([]events.Event, error) {
	e, ok := ev.(events.UserEnteredDungeonEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected event type")
	}

	session.State = models.PlayerStateInDungeon
	session.DungeonEnteredAt = e.Timestamp
	session.CurrentFloor = 0
	session.CurrentFloorStartedAt = e.Timestamp
	session.CurrentFloorFinishedAt = time.Time{}
	session.MonstersKilledOnFloor = 0
	session.Metrics = models.Metrics{
		FloorDurations: make([]time.Duration, 0),
	}

	return []events.Event{
		events.UserEnteredDungeonEvent{BaseEvent: events.BaseEvent{Timestamp: e.Timestamp, UserID: e.UserID}},
	}, nil
}
