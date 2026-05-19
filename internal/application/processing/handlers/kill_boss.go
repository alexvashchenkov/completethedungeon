package handlers

import (
	"fmt"
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
)

type KillBossHandler struct{}

func NewKillBossHandler() *KillBossHandler {
	return &KillBossHandler{}
}

func (h *KillBossHandler) Handle(session *models.PlayerSession, ev events.Event) ([]events.Event, error) {
	e, ok := ev.(events.UserKilledBossEvent)
	if !ok {
		return nil, fmt.Errorf("unexpected event type")
	}

	session.BossKilled = true
	session.BossKilledAt = e.Timestamp
	session.State = models.PlayerStateCompleted

	if !session.CurrentFloorStartedAt.IsZero() {
		session.Metrics.BossDuration = e.Timestamp.Sub(session.CurrentFloorStartedAt)
		session.CurrentFloorFinishedAt = e.Timestamp
	}

	return []events.Event{
		events.UserKilledBossEvent{BaseEvent: events.BaseEvent{Timestamp: e.Timestamp, UserID: e.UserID}},
	}, nil
}
