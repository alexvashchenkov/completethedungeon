package handlers

import (
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
	"testing"
	"time"
)

func TestNextFloorHandler_Handle(t *testing.T) {
	tests := []struct {
		name                    string
		cfg                     *models.Config
		initialFloor            int
		initialMonstersKilled   int
		initialFloorStartedAt   time.Time
		initialFloorFinishedAt  time.Time
		expectSuccess           bool
		expectedOutEventType    string
		expectedNewFloor        int
		expectedMonstersReset   bool
		expectedFloorsCompleted int
	}{
		{
			name:                    "cannot move to next floor - not all monsters killed",
			cfg:                     &models.Config{Floors: 3, Monsters: 5},
			initialFloor:            0,
			initialMonstersKilled:   3,
			initialFloorStartedAt:   time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC),
			initialFloorFinishedAt:  time.Time{},
			expectSuccess:           false,
			expectedOutEventType:    "ImpossibleMove",
			expectedNewFloor:        0,
			expectedMonstersReset:   false,
			expectedFloorsCompleted: 0,
		},
		{
			name:                    "move to next floor successfully",
			cfg:                     &models.Config{Floors: 3, Monsters: 5},
			initialFloor:            0,
			initialMonstersKilled:   5,
			initialFloorStartedAt:   time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC),
			initialFloorFinishedAt:  time.Date(2024, 5, 19, 10, 5, 0, 0, time.UTC),
			expectSuccess:           true,
			expectedOutEventType:    "NextFloor",
			expectedNewFloor:        1,
			expectedMonstersReset:   true,
			expectedFloorsCompleted: 1,
		},
		{
			name:                    "cannot move to next floor - already at last regular floor",
			cfg:                     &models.Config{Floors: 3, Monsters: 5},
			initialFloor:            2,
			initialMonstersKilled:   5,
			initialFloorStartedAt:   time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC),
			initialFloorFinishedAt:  time.Date(2024, 5, 19, 10, 5, 0, 0, time.UTC),
			expectSuccess:           false,
			expectedOutEventType:    "ImpossibleMove",
			expectedNewFloor:        2,
			expectedMonstersReset:   false,
			expectedFloorsCompleted: 0,
		},
		{
			name:                    "nil config",
			cfg:                     nil,
			initialFloor:            0,
			initialMonstersKilled:   5,
			initialFloorStartedAt:   time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC),
			initialFloorFinishedAt:  time.Date(2024, 5, 19, 10, 5, 0, 0, time.UTC),
			expectSuccess:           false,
			expectedOutEventType:    "ImpossibleMove",
			expectedNewFloor:        0,
			expectedMonstersReset:   false,
			expectedFloorsCompleted: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewNextFloorHandler(tt.cfg)

			timestamp := time.Date(2024, 5, 19, 10, 10, 0, 0, time.UTC)
			session := &models.PlayerSession{
				ID:                     1,
				CurrentFloor:           tt.initialFloor,
				MonstersKilledOnFloor:  tt.initialMonstersKilled,
				CurrentFloorStartedAt:  tt.initialFloorStartedAt,
				CurrentFloorFinishedAt: tt.initialFloorFinishedAt,
			}

			event := events.UserWentToNextFloorEvent{
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

			_, isNextFloor := outs[0].(events.UserWentToNextFloorEvent)
			_, isImpossible := outs[0].(events.UserMakesImpossibleMove)

			if tt.expectSuccess && !isNextFloor {
				t.Errorf("expected NextFloor event, got %T", outs[0])
			}

			if !tt.expectSuccess && !isImpossible {
				t.Errorf("expected ImpossibleMove event, got %T", outs[0])
			}

			if session.CurrentFloor != tt.expectedNewFloor {
				t.Errorf("expected floor %d, got %d", tt.expectedNewFloor, session.CurrentFloor)
			}

			if tt.expectedMonstersReset && session.MonstersKilledOnFloor != 0 {
				t.Errorf("expected monsters to be reset to 0, got %d", session.MonstersKilledOnFloor)
			}

			if tt.expectSuccess && session.Metrics.FloorsCompleted != tt.expectedFloorsCompleted {
				t.Errorf("expected %d floors completed, got %d", tt.expectedFloorsCompleted, session.Metrics.FloorsCompleted)
			}
		})
	}
}

func TestNextFloorHandler_InvalidEventType(t *testing.T) {
	handler := NewNextFloorHandler(&models.Config{Floors: 3, Monsters: 5})
	session := &models.PlayerSession{ID: 1}

	wrongEvent := events.UserKilledMonsterEvent{
		BaseEvent: events.BaseEvent{Timestamp: time.Now(), UserID: 1},
	}

	_, err := handler.Handle(session, wrongEvent)

	if err == nil {
		t.Fatal("expected error for invalid event type")
	}
}

func TestNextFloorHandler_FloorDurationTracking(t *testing.T) {
	cfg := &models.Config{Floors: 3, Monsters: 5}
	handler := NewNextFloorHandler(cfg)

	startTime := time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC)
	finishTime := time.Date(2024, 5, 19, 10, 5, 0, 0, time.UTC)

	session := &models.PlayerSession{
		ID:                     1,
		CurrentFloor:           0,
		MonstersKilledOnFloor:  5,
		CurrentFloorStartedAt:  startTime,
		CurrentFloorFinishedAt: finishTime,
	}

	event := events.UserWentToNextFloorEvent{
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
