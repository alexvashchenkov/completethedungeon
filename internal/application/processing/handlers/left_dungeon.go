package handlers

import (
	"fmt"
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
)

type LeftDungeonHandler struct{}

func NewLeftDungeonHandler() *LeftDungeonHandler {
	return &LeftDungeonHandler{}
}

func (h *LeftDungeonHandler) Handle(session *models.PlayerSession, ev events.Event) ([]events.Event, error) {
	e, ok := ev.(events.UserLeftDungeonEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected event type")
	}

	session.DungeonLeftAt = e.Timestamp
	session.DungeonFinishedAt = e.Timestamp
	session.State = models.PlayerStateLeftDungeon

	if !session.DungeonEnteredAt.IsZero() {
		session.Metrics.TotalDungeonDuration = e.Timestamp.Sub(session.DungeonEnteredAt)
	}

	return []events.Event{
		events.UserLeftDungeonEvent{BaseEvent: events.BaseEvent{Timestamp: e.Timestamp, UserID: e.UserID}},
	}, nil
}
