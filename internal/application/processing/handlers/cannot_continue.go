package handlers

import (
	"fmt"
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
)

type CannotContinueHandler struct{}

func NewCannotContinueHandler() *CannotContinueHandler {
	return &CannotContinueHandler{}
}

func (h *CannotContinueHandler) Handle(session *models.PlayerSession, ev events.Event) ([]events.Event, error) {
	e, ok := ev.(events.UserCantContinueDueToReasonEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected event type")
	}

	session.State = models.PlayerStateDisqualified
	session.DungeonFinishedAt = e.Timestamp
	if !session.DungeonEnteredAt.IsZero() {
		session.Metrics.TotalDungeonDuration = e.Timestamp.Sub(session.DungeonEnteredAt)
	}

	return []events.Event{
		events.UserDisqualifiedEvent{BaseEvent: events.BaseEvent{Timestamp: e.Timestamp, UserID: e.UserID}},
	}, nil
}
