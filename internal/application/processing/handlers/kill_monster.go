package handlers

import (
	"fmt"
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
)

type KillMonsterHandler struct{}

func NewKillMonsterHandler() *KillMonsterHandler {
	return &KillMonsterHandler{}
}

func (h *KillMonsterHandler) Handle(session *models.PlayerSession, ev events.Event) ([]events.Event, error) {
	e, ok := ev.(events.UserKilledMonsterEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected event type")
	}

	out := []events.Event{
		events.UserKilledMonsterEvent{BaseEvent: events.BaseEvent{Timestamp: e.Timestamp, UserID: e.UserID}},
	}
	session.MonstersKilledOnFloor += 1

	return out, nil
}
