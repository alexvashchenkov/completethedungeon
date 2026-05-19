package handlers

import (
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
	"testing"
	"time"
)

func TestCannotContinueHandler_Handle(t *testing.T) {
	handler := NewCannotContinueHandler()

	dungeonEnteredAt := time.Date(2024, 5, 19, 9, 0, 0, 0, time.UTC)
	cannotContinueTime := time.Date(2024, 5, 19, 10, 15, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:               1,
		State:            models.PlayerStateInDungeon,
		DungeonEnteredAt: dungeonEnteredAt,
	}

	event := events.UserCantContinueDueToReasonEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: cannotContinueTime,
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

	if session.State != models.PlayerStateDisqualified {
		t.Errorf("expected state PlayerStateDisqualified, got %v", session.State)
	}

	if session.DungeonFinishedAt != cannotContinueTime {
		t.Errorf("expected DungeonFinishedAt to be %v, got %v", cannotContinueTime, session.DungeonFinishedAt)
	}

	expectedDuration := cannotContinueTime.Sub(dungeonEnteredAt)
	if session.Metrics.TotalDungeonDuration != expectedDuration {
		t.Errorf("expected TotalDungeonDuration %v, got %v", expectedDuration, session.Metrics.TotalDungeonDuration)
	}

	_, isEvent := outs[0].(events.UserDisqualifiedEvent)
	if !isEvent {
		t.Errorf("expected UserDisqualifiedEvent, got %T", outs[0])
	}
}

func TestCannotContinueHandler_InvalidEventType(t *testing.T) {
	handler := NewCannotContinueHandler()
	session := &models.PlayerSession{ID: 1}

	wrongEvent := events.UserKilledMonsterEvent{
		BaseEvent: events.BaseEvent{Timestamp: time.Now(), UserID: 1},
	}

	_, err := handler.Handle(session, wrongEvent)

	if err == nil {
		t.Fatal("expected error for invalid event type")
	}
}

func TestCannotContinueHandler_NoDungeonEnteredTime(t *testing.T) {
	handler := NewCannotContinueHandler()

	cannotContinueTime := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:    1,
		State: models.PlayerStateInDungeon,
	}

	event := events.UserCantContinueDueToReasonEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: cannotContinueTime,
			UserID:    1,
		},
	}

	_, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.State != models.PlayerStateDisqualified {
		t.Errorf("expected state PlayerStateDisqualified, got %v", session.State)
	}

	if session.Metrics.TotalDungeonDuration != 0 {
		t.Errorf("expected TotalDungeonDuration to be 0, got %v", session.Metrics.TotalDungeonDuration)
	}

	if session.DungeonFinishedAt != cannotContinueTime {
		t.Errorf("expected DungeonFinishedAt to be %v, got %v", cannotContinueTime, session.DungeonFinishedAt)
	}
}

func TestCannotContinueHandler_DisqualifiesDuringAttempt(t *testing.T) {
	handler := NewCannotContinueHandler()

	dungeonEnteredAt := time.Date(2024, 5, 19, 9, 0, 0, 0, time.UTC)
	cannotContinueTime := time.Date(2024, 5, 19, 9, 45, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:                    1,
		State:                 models.PlayerStateInDungeon,
		DungeonEnteredAt:      dungeonEnteredAt,
		CurrentFloor:          1,
		MonstersKilledOnFloor: 3,
		Metrics: models.Metrics{
			FloorsCompleted: 1,
		},
	}

	event := events.UserCantContinueDueToReasonEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: cannotContinueTime,
			UserID:    1,
		},
	}

	_, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.State != models.PlayerStateDisqualified {
		t.Errorf("expected state PlayerStateDisqualified, got %v", session.State)
	}

	expectedDuration := cannotContinueTime.Sub(dungeonEnteredAt)
	if session.Metrics.TotalDungeonDuration != expectedDuration {
		t.Errorf("expected TotalDungeonDuration %v, got %v", expectedDuration, session.Metrics.TotalDungeonDuration)
	}

	if session.CurrentFloor != 1 {
		t.Errorf("expected CurrentFloor to remain 1, got %d", session.CurrentFloor)
	}
}
