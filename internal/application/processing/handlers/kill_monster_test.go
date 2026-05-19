package handlers

import (
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
	"testing"
	"time"
)

func TestKillMonsterHandler_Handle(t *testing.T) {
	tests := []struct {
		name                   string
		cfg                    *models.Config
		initialMonstersKilled  int
		expectedMonstersKilled int
		expectFloorFinished    bool
	}{
		{
			name:                   "kill monster within limit",
			cfg:                    &models.Config{Floors: 2, Monsters: 3},
			initialMonstersKilled:  0,
			expectedMonstersKilled: 1,
			expectFloorFinished:    false,
		},
		{
			name:                   "kill last monster on floor",
			cfg:                    &models.Config{Floors: 2, Monsters: 3},
			initialMonstersKilled:  2,
			expectedMonstersKilled: 3,
			expectFloorFinished:    true,
		},
		{
			name:                   "nil config",
			cfg:                    nil,
			initialMonstersKilled:  0,
			expectedMonstersKilled: 1,
			expectFloorFinished:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewKillMonsterHandler(tt.cfg)

			timestamp := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
			session := &models.PlayerSession{
				ID:                    1,
				CurrentFloor:          0,
				MonstersKilledOnFloor: tt.initialMonstersKilled,
			}

			event := events.UserKilledMonsterEvent{
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

			if session.MonstersKilledOnFloor != tt.expectedMonstersKilled {
				t.Errorf("expected %d monsters killed, got %d", tt.expectedMonstersKilled, session.MonstersKilledOnFloor)
			}

			if tt.expectFloorFinished && session.CurrentFloorFinishedAt.IsZero() {
				t.Error("expected CurrentFloorFinishedAt to be set, but it's zero")
			}

			if !tt.expectFloorFinished && !session.CurrentFloorFinishedAt.IsZero() {
				t.Error("expected CurrentFloorFinishedAt to be zero, but it's set")
			}
		})
	}
}

func TestKillMonsterHandler_InvalidEventType(t *testing.T) {
	handler := NewKillMonsterHandler(&models.Config{Floors: 2, Monsters: 3})
	session := &models.PlayerSession{ID: 1}

	wrongEvent := events.UserWentToNextFloorEvent{
		BaseEvent: events.BaseEvent{Timestamp: time.Now(), UserID: 1},
	}

	_, err := handler.Handle(session, wrongEvent)

	if err == nil {
		t.Fatal("expected error for invalid event type")
	}
}

func TestKillMonsterHandler_CannotExceedMaxMonsters(t *testing.T) {
	cfg := &models.Config{Floors: 2, Monsters: 3}
	handler := NewKillMonsterHandler(cfg)

	timestamp := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	session := &models.PlayerSession{
		ID:                    1,
		CurrentFloor:          0,
		MonstersKilledOnFloor: 3,
	}

	event := events.UserKilledMonsterEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: timestamp,
			UserID:    1,
		},
	}

	outs, err := handler.Handle(session, event)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.MonstersKilledOnFloor != 4 {
		t.Errorf("expected monsters killed to be incremented to 4, got %d", session.MonstersKilledOnFloor)
	}

	if len(outs) != 1 {
		t.Fatalf("expected 1 output event, got %d", len(outs))
	}
}
