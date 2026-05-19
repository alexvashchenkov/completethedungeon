package handlers

import (
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
	"testing"
	"time"
)

func TestRestoredHealthHandler_IncreaseHP(t *testing.T) {
	handler := NewRestoredHealthHandler()

	timestamp := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:    1,
		HP:    50,
		State: models.PlayerStateInDungeon,
	}

	event := events.UserRestoredHealthEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: timestamp,
			UserID:    1,
		},
		Amount: 20,
	}

	outs, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(outs) != 1 {
		t.Fatalf("expected 1 output event, got %d", len(outs))
	}

	if session.HP != 70 {
		t.Errorf("expected HP to be 70, got %d", session.HP)
	}

	outEvent, isEvent := outs[0].(events.UserRestoredHealthEvent)
	if !isEvent {
		t.Errorf("expected UserRestoredHealthEvent, got %T", outs[0])
	}

	if outEvent.Amount != 20 {
		t.Errorf("expected event amount to be 20, got %d", outEvent.Amount)
	}
}

func TestRestoredHealthHandler_CapAtMax(t *testing.T) {
	handler := NewRestoredHealthHandler()

	timestamp := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:    1,
		HP:    90,
		State: models.PlayerStateInDungeon,
	}

	event := events.UserRestoredHealthEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: timestamp,
			UserID:    1,
		},
		Amount: 50,
	}

	outs, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(outs) != 1 {
		t.Fatalf("expected 1 output event, got %d", len(outs))
	}

	if session.HP != 100 {
		t.Errorf("expected HP to be capped at 100, got %d", session.HP)
	}
}

func TestRestoredHealthHandler_RestoreToExactMax(t *testing.T) {
	handler := NewRestoredHealthHandler()

	timestamp := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:    1,
		HP:    85,
		State: models.PlayerStateInDungeon,
	}

	event := events.UserRestoredHealthEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: timestamp,
			UserID:    1,
		},
		Amount: 15,
	}

	outs, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(outs) != 1 {
		t.Fatalf("expected 1 output event, got %d", len(outs))
	}

	if session.HP != 100 {
		t.Errorf("expected HP to be exactly 100, got %d", session.HP)
	}
}

func TestRestoredHealthHandler_InvalidEventType(t *testing.T) {
	handler := NewRestoredHealthHandler()
	session := &models.PlayerSession{ID: 1}

	wrongEvent := events.UserKilledMonsterEvent{
		BaseEvent: events.BaseEvent{Timestamp: time.Now(), UserID: 1},
	}

	_, err := handler.Handle(session, wrongEvent)

	if err == nil {
		t.Fatal("expected error for invalid event type")
	}
}

func TestRestoredHealthHandler_SmallRestoration(t *testing.T) {
	handler := NewRestoredHealthHandler()

	timestamp := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:    1,
		HP:    10,
		State: models.PlayerStateInDungeon,
	}

	event := events.UserRestoredHealthEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: timestamp,
			UserID:    1,
		},
		Amount: 1,
	}

	outs, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(outs) != 1 {
		t.Fatalf("expected 1 output event, got %d", len(outs))
	}

	if session.HP != 11 {
		t.Errorf("expected HP to be 11, got %d", session.HP)
	}
}

func TestRestoredHealthHandler_MaxRestoration(t *testing.T) {
	handler := NewRestoredHealthHandler()

	timestamp := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:    1,
		HP:    1,
		State: models.PlayerStateInDungeon,
	}

	event := events.UserRestoredHealthEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: timestamp,
			UserID:    1,
		},
		Amount: 99,
	}

	_, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.HP != 100 {
		t.Errorf("expected HP to be 100, got %d", session.HP)
	}
}

func TestRestoredHealthHandler_MultipleRestorations(t *testing.T) {
	handler := NewRestoredHealthHandler()

	session := &models.PlayerSession{
		ID:    1,
		HP:    50,
		State: models.PlayerStateInDungeon,
	}

	event1 := events.UserRestoredHealthEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC),
			UserID:    1,
		},
		Amount: 20,
	}

	_, err := handler.Handle(session, event1)
	if err != nil {
		t.Fatalf("unexpected error in first restoration: %v", err)
	}

	if session.HP != 70 {
		t.Errorf("expected HP to be 70 after first restoration, got %d", session.HP)
	}

	event2 := events.UserRestoredHealthEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: time.Date(2024, 5, 19, 10, 5, 0, 0, time.UTC),
			UserID:    1,
		},
		Amount: 35,
	}

	_, err = handler.Handle(session, event2)
	if err != nil {
		t.Fatalf("unexpected error in second restoration: %v", err)
	}

	if session.HP != 100 {
		t.Errorf("expected HP to be capped at 100 after second restoration, got %d", session.HP)
	}
}
