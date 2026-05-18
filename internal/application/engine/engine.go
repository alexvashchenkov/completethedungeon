package engine

import (
	"impulse/internal/application/processing"
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
	"impulse/internal/domain/storage"
)

type Engine struct {
	config      models.Config
	sessionRepo storage.SessionStore
	processor   processing.Processor
}

func NewEngine(sessionRepo storage.SessionStore) *Engine {
	return &Engine{
		sessionRepo: sessionRepo,
	}
}

func (e *Engine) Process(event events.Event) ([]events.Event, error) {
	var session *models.PlayerSession

	if event.ID() == events.UserRegisteredEventID {
		session, _ = e.sessionRepo.Create(event.GetUserID())
		return []events.Event{
			events.UserRegisteredEvent{
				BaseEvent: events.BaseEvent{
					Timestamp: event.GetTimestamp(),
					UserID:    event.GetUserID(),
				},
			},
		}, nil
	}

	session, err := e.sessionRepo.Get(event.GetUserID())
	if err != nil {
		return nil, err
	}

	switch session.State {
	case models.PlayerStateDead, models.PlayerStateDisqualified,
		models.PlayerStateCompleted, models.PlayerStateLeftDungeon:
		return []events.Event{
			events.UserMakesImpossibleMove{
				BaseEvent: events.BaseEvent{
					Timestamp: event.GetTimestamp(),
					UserID:    event.GetUserID(),
				},
			},
		}, nil
	}

	return e.processor.Process(session, event)
}
