package handlers

import (
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
	"testing"
	"time"
)

func TestLeftDungeonHandler_Handle(t *testing.T) {
	handler := NewLeftDungeonHandler()

	dungeonEnteredAt := time.Date(2024, 5, 19, 9, 0, 0, 0, time.UTC)
	leftTime := time.Date(2024, 5, 19, 10, 30, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:               1,
		State:            models.PlayerStateInDungeon,
		DungeonEnteredAt: dungeonEnteredAt,
	}

	event := events.UserLeftDungeonEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: leftTime,
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

	if session.State != models.PlayerStateLeftDungeon {
		t.Errorf("expected state PlayerStateLeftDungeon, got %v", session.State)
	}

	if session.DungeonLeftAt != leftTime {
		t.Errorf("expected DungeonLeftAt to be %v, got %v", leftTime, session.DungeonLeftAt)
	}

	if session.DungeonFinishedAt != leftTime {
		t.Errorf("expected DungeonFinishedAt to be %v, got %v", leftTime, session.DungeonFinishedAt)
	}

	expectedDuration := leftTime.Sub(dungeonEnteredAt)
	if session.Metrics.TotalDungeonDuration != expectedDuration {
		t.Errorf("expected TotalDungeonDuration %v, got %v", expectedDuration, session.Metrics.TotalDungeonDuration)
	}

	_, isEvent := outs[0].(events.UserLeftDungeonEvent)
	if !isEvent {
		t.Errorf("expected UserLeftDungeonEvent, got %T", outs[0])
	}
}

func TestLeftDungeonHandler_InvalidEventType(t *testing.T) {
	handler := NewLeftDungeonHandler()
	session := &models.PlayerSession{ID: 1}

	wrongEvent := events.UserKilledMonsterEvent{
		BaseEvent: events.BaseEvent{Timestamp: time.Now(), UserID: 1},
	}

	_, err := handler.Handle(session, wrongEvent)

	if err == nil {
		t.Fatal("expected error for invalid event type")
	}
}

func TestLeftDungeonHandler_NoDungeonEnteredTime(t *testing.T) {
	handler := NewLeftDungeonHandler()

	leftTime := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:    1,
		State: models.PlayerStateInDungeon,
	}

	event := events.UserLeftDungeonEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: leftTime,
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

	if session.State != models.PlayerStateLeftDungeon {
		t.Errorf("expected state PlayerStateLeftDungeon, got %v", session.State)
	}

	if session.Metrics.TotalDungeonDuration != 0 {
		t.Errorf("expected TotalDungeonDuration to be 0, got %v", session.Metrics.TotalDungeonDuration)
	}
}

func TestLeftDungeonHandler_AfterCompletion(t *testing.T) {
	handler := NewLeftDungeonHandler()

	dungeonEnteredAt := time.Date(2024, 5, 19, 9, 0, 0, 0, time.UTC)
	leftTime := time.Date(2024, 5, 19, 10, 30, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:               1,
		State:            models.PlayerStateCompleted,
		DungeonEnteredAt: dungeonEnteredAt,
	}

	event := events.UserLeftDungeonEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: leftTime,
			UserID:    1,
		},
	}

	_, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.State != models.PlayerStateLeftDungeon {
		t.Errorf("expected state PlayerStateLeftDungeon, got %v", session.State)
	}

	expectedDuration := leftTime.Sub(dungeonEnteredAt)
	if session.Metrics.TotalDungeonDuration != expectedDuration {
		t.Errorf("expected TotalDungeonDuration %v, got %v", expectedDuration, session.Metrics.TotalDungeonDuration)
	}
}
