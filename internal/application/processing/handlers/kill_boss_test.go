package handlers

import (
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
	"testing"
	"time"
)

func TestKillBossHandler_Handle(t *testing.T) {
	handler := NewKillBossHandler()

	bossStartTime := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	bossKillTime := time.Date(2024, 5, 19, 10, 15, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:                    1,
		State:                 models.PlayerStateInDungeon,
		CurrentFloor:          3,
		CurrentFloorStartedAt: bossStartTime,
	}

	event := events.UserKilledBossEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: bossKillTime,
			UserID:    1,
		},
	}

	outs, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(outs) != 1 {
		t.Fatalf("expected 1 output event, got %d", len(outs))
	}

	if !session.BossKilled {
		t.Error("expected BossKilled to be true")
	}

	if session.BossKilledAt != bossKillTime {
		t.Errorf("expected BossKilledAt to be %v, got %v", bossKillTime, session.BossKilledAt)
	}

	if session.State != models.PlayerStateCompleted {
		t.Errorf("expected state PlayerStateCompleted, got %v", session.State)
	}

	expectedDuration := bossKillTime.Sub(bossStartTime)
	if session.Metrics.BossDuration != expectedDuration {
		t.Errorf("expected BossDuration %v, got %v", expectedDuration, session.Metrics.BossDuration)
	}

	if session.CurrentFloorFinishedAt != bossKillTime {
		t.Errorf("expected CurrentFloorFinishedAt to be %v, got %v", bossKillTime, session.CurrentFloorFinishedAt)
	}

	_, isEvent := outs[0].(events.UserKilledBossEvent)
	if !isEvent {
		t.Errorf("expected UserKilledBossEvent, got %T", outs[0])
	}
}

func TestKillBossHandler_InvalidEventType(t *testing.T) {
	handler := NewKillBossHandler()
	session := &models.PlayerSession{ID: 1}

	wrongEvent := events.UserKilledMonsterEvent{
		BaseEvent: events.BaseEvent{Timestamp: time.Now(), UserID: 1},
	}

	_, err := handler.Handle(session, wrongEvent)

	if err == nil {
		t.Fatal("expected error for invalid event type")
	}
}

func TestKillBossHandler_NoFloorStartTime(t *testing.T) {
	handler := NewKillBossHandler()

	killTime := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:                    1,
		State:                 models.PlayerStateInDungeon,
		CurrentFloor:          3,
		CurrentFloorStartedAt: time.Time{},
	}

	event := events.UserKilledBossEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: killTime,
			UserID:    1,
		},
	}

	outs, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(outs) != 1 {
		t.Fatalf("expected 1 output event, got %d", len(outs))
	}

	if session.Metrics.BossDuration != 0 {
		t.Errorf("expected BossDuration to be 0, got %v", session.Metrics.BossDuration)
	}

	if !session.BossKilled {
		t.Error("expected BossKilled to be true")
	}

	if session.State != models.PlayerStateCompleted {
		t.Errorf("expected state PlayerStateCompleted, got %v", session.State)
	}
}

func TestKillBossHandler_StateTransition(t *testing.T) {
	handler := NewKillBossHandler()

	session := &models.PlayerSession{
		ID:                    1,
		State:                 models.PlayerStateInDungeon,
		CurrentFloor:          3,
		CurrentFloorStartedAt: time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC),
	}

	event := events.UserKilledBossEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: time.Date(2024, 5, 19, 10, 10, 0, 0, time.UTC),
			UserID:    1,
		},
	}

	_, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.State != models.PlayerStateCompleted {
		t.Errorf("expected state transition to PlayerStateCompleted, got %v", session.State)
	}
}
