package handlers

import (
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
	"testing"
	"time"
)

func TestPreviousFloorHandler_Handle(t *testing.T) {
	cfg := &models.Config{Floors: 3, Monsters: 5}
	handler := NewPreviousFloorHandler(cfg)

	timestamp := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:                    1,
		CurrentFloor:          2,
		MonstersKilledOnFloor: 3,
	}

	event := events.UserWentToPreviousFloorEvent{
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

	if session.CurrentFloor != 1 {
		t.Errorf("expected current floor to be 1, got %d", session.CurrentFloor)
	}

	if session.MonstersKilledOnFloor != 0 {
		t.Errorf("expected MonstersKilledOnFloor to be reset to 0, got %d", session.MonstersKilledOnFloor)
	}

	_, isEvent := outs[0].(events.UserWentToPreviousFloorEvent)
	if !isEvent {
		t.Errorf("expected UserWentToPreviousFloorEvent, got %T", outs[0])
	}
}

func TestPreviousFloorHandler_CannotGoBeforeFloorZero(t *testing.T) {
	cfg := &models.Config{Floors: 3, Monsters: 5}
	handler := NewPreviousFloorHandler(cfg)

	timestamp := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:                    1,
		CurrentFloor:          0,
		MonstersKilledOnFloor: 2,
	}

	event := events.UserWentToPreviousFloorEvent{
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

	if session.CurrentFloor != 0 {
		t.Errorf("expected current floor to remain 0, got %d", session.CurrentFloor)
	}

	if session.MonstersKilledOnFloor != 2 {
		t.Errorf("expected MonstersKilledOnFloor to remain 2, got %d", session.MonstersKilledOnFloor)
	}

	_, isImpossible := outs[0].(events.UserMakesImpossibleMove)
	if !isImpossible {
		t.Errorf("expected UserMakesImpossibleMove, got %T", outs[0])
	}
}

func TestPreviousFloorHandler_InvalidEventType(t *testing.T) {
	handler := NewPreviousFloorHandler(&models.Config{Floors: 3, Monsters: 5})
	session := &models.PlayerSession{ID: 1}

	wrongEvent := events.UserKilledMonsterEvent{
		BaseEvent: events.BaseEvent{Timestamp: time.Now(), UserID: 1},
	}

	_, err := handler.Handle(session, wrongEvent)

	if err == nil {
		t.Fatal("expected error for invalid event type")
	}
}

func TestPreviousFloorHandler_MultipleBackwardMoves(t *testing.T) {
	cfg := &models.Config{Floors: 3, Monsters: 5}
	handler := NewPreviousFloorHandler(cfg)

	session := &models.PlayerSession{
		ID:                    1,
		CurrentFloor:          3,
		MonstersKilledOnFloor: 5,
	}

	timestamp1 := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	event1 := events.UserWentToPreviousFloorEvent{
		BaseEvent: events.BaseEvent{Timestamp: timestamp1, UserID: 1},
	}

	outs1, err := handler.Handle(session, event1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, isEvent := outs1[0].(events.UserWentToPreviousFloorEvent)
	if !isEvent {
		t.Errorf("expected successful move event, got %T", outs1[0])
	}

	if session.CurrentFloor != 2 {
		t.Errorf("expected floor 2 after first move, got %d", session.CurrentFloor)
	}

	timestamp2 := time.Date(2024, 5, 19, 10, 5, 0, 0, time.UTC)
	event2 := events.UserWentToPreviousFloorEvent{
		BaseEvent: events.BaseEvent{Timestamp: timestamp2, UserID: 1},
	}

	outs2, err := handler.Handle(session, event2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, isEvent = outs2[0].(events.UserWentToPreviousFloorEvent)
	if !isEvent {
		t.Errorf("expected successful move event, got %T", outs2[0])
	}

	if session.CurrentFloor != 1 {
		t.Errorf("expected floor 1 after second move, got %d", session.CurrentFloor)
	}

	timestamp3 := time.Date(2024, 5, 19, 10, 10, 0, 0, time.UTC)
	event3 := events.UserWentToPreviousFloorEvent{
		BaseEvent: events.BaseEvent{Timestamp: timestamp3, UserID: 1},
	}

	outs3, err := handler.Handle(session, event3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, isEvent = outs3[0].(events.UserWentToPreviousFloorEvent)
	if !isEvent {
		t.Errorf("expected successful move event, got %T", outs3[0])
	}

	if session.CurrentFloor != 0 {
		t.Errorf("expected floor to be 0 after third move, got %d", session.CurrentFloor)
	}

	timestamp4 := time.Date(2024, 5, 19, 10, 15, 0, 0, time.UTC)
	event4 := events.UserWentToPreviousFloorEvent{
		BaseEvent: events.BaseEvent{Timestamp: timestamp4, UserID: 1},
	}

	outs4, err := handler.Handle(session, event4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, isImpossible := outs4[0].(events.UserMakesImpossibleMove)
	if !isImpossible {
		t.Errorf("expected impossible move event, got %T", outs4[0])
	}

	if session.CurrentFloor != 0 {
		t.Errorf("expected floor to remain 0, got %d", session.CurrentFloor)
	}
}

func TestPreviousFloorHandler_ResetsMonsterCounter(t *testing.T) {
	cfg := &models.Config{Floors: 3, Monsters: 5}
	handler := NewPreviousFloorHandler(cfg)

	session := &models.PlayerSession{
		ID:                    1,
		CurrentFloor:          2,
		MonstersKilledOnFloor: 5,
	}

	event := events.UserWentToPreviousFloorEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC),
			UserID:    1,
		},
	}

	_, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.MonstersKilledOnFloor != 0 {
		t.Errorf("expected MonstersKilledOnFloor to be reset to 0, got %d", session.MonstersKilledOnFloor)
	}
}
