package handlers

import (
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
	"testing"
	"time"
)

func TestEnteredBossFloorHandler_Handle(t *testing.T) {
	tests := []struct {
		name                   string
		cfg                    *models.Config
		initialFloor           int
		initialFloorStartedAt  time.Time
		initialFloorFinishedAt time.Time
		expectSuccess          bool
		expectedOutEventType   string
		expectedNewFloor       int
		expectedMonsterReset   bool
	}{
		{
			name:                   "enter boss floor successfully",
			cfg:                    &models.Config{Floors: 3, Monsters: 5},
			initialFloor:           2,
			initialFloorStartedAt:  time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC),
			initialFloorFinishedAt: time.Date(2024, 5, 19, 10, 5, 0, 0, time.UTC),
			expectSuccess:          true,
			expectedOutEventType:   "BossFloor",
			expectedNewFloor:       3,
			expectedMonsterReset:   true,
		},
		{
			name:                   "cannot enter boss floor - already at or beyond boss floor",
			cfg:                    &models.Config{Floors: 3, Monsters: 5},
			initialFloor:           3,
			initialFloorStartedAt:  time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC),
			initialFloorFinishedAt: time.Date(2024, 5, 19, 10, 5, 0, 0, time.UTC),
			expectSuccess:          false,
			expectedOutEventType:   "ImpossibleMove",
			expectedNewFloor:       3,
			expectedMonsterReset:   false,
		},
		{
			name:                   "clear floor duration when entering boss floor",
			cfg:                    &models.Config{Floors: 3, Monsters: 5},
			initialFloor:           1,
			initialFloorStartedAt:  time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC),
			initialFloorFinishedAt: time.Date(2024, 5, 19, 10, 5, 0, 0, time.UTC),
			expectSuccess:          true,
			expectedOutEventType:   "BossFloor",
			expectedNewFloor:       3,
			expectedMonsterReset:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewEnteredBossFloorHandler(tt.cfg)

			timestamp := time.Date(2024, 5, 19, 10, 10, 0, 0, time.UTC)
			session := &models.PlayerSession{
				ID:                     1,
				CurrentFloor:           tt.initialFloor,
				MonstersKilledOnFloor:  2,
				CurrentFloorStartedAt:  tt.initialFloorStartedAt,
				CurrentFloorFinishedAt: tt.initialFloorFinishedAt,
			}

			event := events.UserEnteredBossFloorEvent{
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

			_, isBossFloor := outs[0].(events.UserEnteredBossFloorEvent)
			_, isImpossible := outs[0].(events.UserMakesImpossibleMove)

			if tt.expectSuccess && !isBossFloor {
				t.Errorf("expected BossFloor event, got %T", outs[0])
			}

			if !tt.expectSuccess && !isImpossible {
				t.Errorf("expected ImpossibleMove event, got %T", outs[0])
			}

			if session.CurrentFloor != tt.expectedNewFloor {
				t.Errorf("expected floor %d, got %d", tt.expectedNewFloor, session.CurrentFloor)
			}

			if tt.expectSuccess && session.MonstersKilledOnFloor != 0 {
				t.Errorf("expected monsters to be reset to 0, got %d", session.MonstersKilledOnFloor)
			}

			if tt.expectSuccess && !session.CurrentFloorFinishedAt.IsZero() {
				t.Error("expected CurrentFloorFinishedAt to be cleared (zero) after entering boss floor")
			}
		})
	}
}

func TestEnteredBossFloorHandler_InvalidEventType(t *testing.T) {
	handler := NewEnteredBossFloorHandler(&models.Config{Floors: 3, Monsters: 5})
	session := &models.PlayerSession{ID: 1}

	wrongEvent := events.UserKilledMonsterEvent{
		BaseEvent: events.BaseEvent{Timestamp: time.Now(), UserID: 1},
	}

	_, err := handler.Handle(session, wrongEvent)

	if err == nil {
		t.Fatal("expected error for invalid event type")
	}
}

func TestEnteredBossFloorHandler_FloorDurationTracking(t *testing.T) {
	cfg := &models.Config{Floors: 3, Monsters: 5}
	handler := NewEnteredBossFloorHandler(cfg)

	startTime := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	finishTime := time.Date(2024, 5, 19, 10, 5, 0, 0, time.UTC)

	session := &models.PlayerSession{
		ID:                     1,
		CurrentFloor:           2,
		MonstersKilledOnFloor:  3,
		CurrentFloorStartedAt:  startTime,
		CurrentFloorFinishedAt: finishTime,
	}

	event := events.UserEnteredBossFloorEvent{
		BaseEvent: events.BaseEvent{
			Timestamp: time.Date(2024, 5, 19, 10, 10, 0, 0, time.UTC),
			UserID:    1,
		},
	}

	_, err := handler.Handle(session, event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(session.Metrics.FloorDurations) != 1 {
		t.Fatalf("expected 1 floor duration, got %d", len(session.Metrics.FloorDurations))
	}

	expectedDuration := finishTime.Sub(startTime)
	if session.Metrics.FloorDurations[0] != expectedDuration {
		t.Errorf("expected duration %v, got %v", expectedDuration, session.Metrics.FloorDurations[0])
	}
}
