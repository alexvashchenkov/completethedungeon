package handlers

import (
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
	"testing"
	"time"
)

func TestDamageHandler_ReduceHP(t *testing.T) {
	handler := NewDamageHandler()

	timestamp := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:               1,
		HP:               100,
		State:            models.PlayerStateInDungeon,
		DungeonEnteredAt: time.Date(2024, 5, 19, 9, 0, 0, 0, time.UTC),
	}

	event := events.UserReceivedDamageEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: timestamp,
			UserID:    1,
		},
		Amount: 10,
	}

	outs, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(outs) != 1 {
		t.Fatalf("expected 1 output event, got %d", len(outs))
	}

	if session.HP != 90 {
		t.Errorf("expected HP to be 90, got %d", session.HP)
	}

	if session.State != models.PlayerStateInDungeon {
		t.Errorf("expected state to remain PlayerStateInDungeon, got %v", session.State)
	}

	_, isDamageEvent := outs[0].(events.UserReceivedDamageEvent)
	if !isDamageEvent {
		t.Errorf("expected UserReceivedDamageEvent, got %T", outs[0])
	}
}

func TestDamageHandler_KillPlayer(t *testing.T) {
	handler := NewDamageHandler()

	dungeonEnteredAt := time.Date(2024, 5, 19, 9, 0, 0, 0, time.UTC)
	damageTime := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:               1,
		HP:               20,
		State:            models.PlayerStateInDungeon,
		DungeonEnteredAt: dungeonEnteredAt,
	}

	event := events.UserReceivedDamageEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: damageTime,
			UserID:    1,
		},
		Amount: 25,
	}

	outs, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(outs) != 2 {
		t.Fatalf("expected 2 output events (damage + death), got %d", len(outs))
	}

	if session.HP != 0 {
		t.Errorf("expected HP to be 0, got %d", session.HP)
	}

	if session.State != models.PlayerStateDead {
		t.Errorf("expected state PlayerStateDead, got %v", session.State)
	}

	if session.DungeonFinishedAt != damageTime {
		t.Errorf("expected DungeonFinishedAt to be %v, got %v", damageTime, session.DungeonFinishedAt)
	}

	expectedDuration := damageTime.Sub(dungeonEnteredAt)
	if session.Metrics.TotalDungeonDuration != expectedDuration {
		t.Errorf("expected duration %v, got %v", expectedDuration, session.Metrics.TotalDungeonDuration)
	}

	_, isDamageEvent := outs[0].(events.UserReceivedDamageEvent)
	if !isDamageEvent {
		t.Errorf("expected first event to be UserReceivedDamageEvent, got %T", outs[0])
	}

	_, isDeadEvent := outs[1].(events.UserDiedEvent)
	if !isDeadEvent {
		t.Errorf("expected second event to be UserDiedEvent, got %T", outs[1])
	}
}

func TestDamageHandler_ExactDeathDamage(t *testing.T) {
	handler := NewDamageHandler()

	dungeonEnteredAt := time.Date(2024, 5, 19, 9, 0, 0, 0, time.UTC)
	damageTime := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:               1,
		HP:               30,
		State:            models.PlayerStateInDungeon,
		DungeonEnteredAt: dungeonEnteredAt,
	}

	event := events.UserReceivedDamageEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: damageTime,
			UserID:    1,
		},
		Amount: 30,
	}

	outs, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(outs) != 2 {
		t.Fatalf("expected 2 output events, got %d", len(outs))
	}

	if session.HP != 0 {
		t.Errorf("expected HP to be 0, got %d", session.HP)
	}

	if session.State != models.PlayerStateDead {
		t.Errorf("expected state PlayerStateDead, got %v", session.State)
	}
}

func TestDamageHandler_InvalidEventType(t *testing.T) {
	handler := NewDamageHandler()
	session := &models.PlayerSession{ID: 1}

	wrongEvent := events.UserKilledMonsterEvent{
		BaseEvent: events.BaseEvent{Timestamp: time.Now(), UserID: 1},
	}

	_, err := handler.Handle(session, wrongEvent)

	if err == nil {
		t.Fatal("expected error for invalid event type")
	}
}

func TestDamageHandler_NoDungeonEnteredTime(t *testing.T) {
	handler := NewDamageHandler()

	session := &models.PlayerSession{
		ID:    1,
		HP:    50,
		State: models.PlayerStateInDungeon,
	}

	event := events.UserReceivedDamageEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC),
			UserID:    1,
		},
		Amount: 60,
	}

	outs, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(outs) != 2 {
		t.Fatalf("expected 2 output events, got %d", len(outs))
	}

	if session.HP != 0 {
		t.Errorf("expected HP to be 0, got %d", session.HP)
	}

	if session.State != models.PlayerStateDead {
		t.Errorf("expected state PlayerStateDead, got %v", session.State)
	}

	if session.Metrics.TotalDungeonDuration != 0 {
		t.Errorf("expected TotalDungeonDuration to be 0, got %v", session.Metrics.TotalDungeonDuration)
	}
}
