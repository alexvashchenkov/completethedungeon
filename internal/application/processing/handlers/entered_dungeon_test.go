package handlers

import (
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
	"testing"
	"time"
)

func TestEnteredDungeonHandler_Handle(t *testing.T) {
	cfg := &models.Config{Floors: 3, Monsters: 5}
	handler := NewEnteredDungeonHandler(cfg)

	timestamp := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:    1,
		HP:    100,
		State: models.PlayerStateRegistered,
	}

	event := events.UserEnteredDungeonEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: timestamp,
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

	if session.State != models.PlayerStateInDungeon {
		t.Errorf("expected state PlayerStateInDungeon, got %v", session.State)
	}

	if session.DungeonEnteredAt != timestamp {
		t.Errorf("expected DungeonEnteredAt to be set to %v, got %v", timestamp, session.DungeonEnteredAt)
	}

	if session.CurrentFloor != 0 {
		t.Errorf("expected initial floor 0, got %d", session.CurrentFloor)
	}

	if session.MonstersKilledOnFloor != 0 {
		t.Errorf("expected MonstersKilledOnFloor to be 0, got %d", session.MonstersKilledOnFloor)
	}

	if session.CurrentFloorStartedAt != timestamp {
		t.Errorf("expected CurrentFloorStartedAt to be %v, got %v", timestamp, session.CurrentFloorStartedAt)
	}

	if !session.CurrentFloorFinishedAt.IsZero() {
		t.Error("expected CurrentFloorFinishedAt to be zero")
	}

	if session.Metrics.FloorDurations == nil {
		t.Error("expected FloorDurations to be initialized")
	}
	if len(session.Metrics.FloorDurations) != 0 {
		t.Errorf("expected empty FloorDurations, got %d", len(session.Metrics.FloorDurations))
	}
}

func TestEnteredDungeonHandler_InvalidEventType(t *testing.T) {
	handler := NewEnteredDungeonHandler(&models.Config{Floors: 3, Monsters: 5})
	session := &models.PlayerSession{ID: 1}

	wrongEvent := events.UserKilledMonsterEvent{
		BaseEvent: events.BaseEvent{Timestamp: time.Now(), UserID: 1},
	}

	_, err := handler.Handle(session, wrongEvent)

	if err == nil {
		t.Fatal("expected error for invalid event type")
	}
}

func TestEnteredDungeonHandler_ResetsPreviousState(t *testing.T) {
	cfg := &models.Config{Floors: 3, Monsters: 5}
	handler := NewEnteredDungeonHandler(cfg)

	timestamp := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:                     1,
		State:                  models.PlayerStateLeftDungeon,
		CurrentFloor:           2,
		MonstersKilledOnFloor:  3,
		CurrentFloorStartedAt:  time.Date(2024, 5, 19, 9, 0, 0, 0, time.UTC),
		CurrentFloorFinishedAt: time.Date(2024, 5, 19, 9, 5, 0, 0, time.UTC),
		Metrics: models.Metrics{
			FloorsCompleted: 2,
			FloorDurations:  []time.Duration{5 * time.Minute, 3 * time.Minute},
		},
	}

	event := events.UserEnteredDungeonEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: timestamp,
			UserID:    1,
		},
	}

	_, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.State != models.PlayerStateInDungeon {
		t.Errorf("expected state PlayerStateInDungeon, got %v", session.State)
	}

	if session.CurrentFloor != 0 {
		t.Errorf("expected floor to be reset to 0, got %d", session.CurrentFloor)
	}

	if session.MonstersKilledOnFloor != 0 {
		t.Errorf("expected MonstersKilledOnFloor to be reset to 0, got %d", session.MonstersKilledOnFloor)
	}

	if len(session.Metrics.FloorDurations) != 0 {
		t.Errorf("expected FloorDurations to be reset to empty, got %d durations", len(session.Metrics.FloorDurations))
	}
}
