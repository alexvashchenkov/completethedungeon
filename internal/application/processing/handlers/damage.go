package handlers

import (
	"fmt"
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
)

type DamageHandler struct {
}

func NewDamageHandler() *DamageHandler { return &DamageHandler{} }

func (h *DamageHandler) Handle(session *models.PlayerSession, ev events.Event) ([]events.Event, error) {
	e, ok := ev.(events.UserReceivedDamageEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected event type")
	}

	out := []events.Event{
		events.UserReceivedDamageEvent{BaseEvent: events.BaseEvent{Timestamp: e.Timestamp, UserID: e.UserID}, Amount: e.Amount},
	}

	session.HP -= e.Amount
	if session.HP <= 0 {
		session.HP = 0
		session.State = models.PlayerStateDead
		session.DungeonFinishedAt = e.Timestamp
		if !session.DungeonEnteredAt.IsZero() {
			session.Metrics.TotalDungeonDuration = e.Timestamp.Sub(session.DungeonEnteredAt)
		}
		out = append(out, events.UserDiedEvent{BaseEvent: events.BaseEvent{Timestamp: e.Timestamp, UserID: e.UserID}})
	}
	return out, nil
}
