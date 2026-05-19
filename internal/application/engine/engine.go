package engine

import (
	"impulse/internal/application/processing"
	"impulse/internal/application/report"
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
	"impulse/internal/domain/storage"
)

type Engine struct {
	sessionRepo     storage.SessionStore
	processor       *processing.Processor
	reportBuilder   report.Builder
	reportFormatter report.Formatter
}

func NewEngine(sessionRepo storage.SessionStore, processor *processing.Processor, reportBuilder report.Builder, reportFormatter report.Formatter) *Engine {
	return &Engine{
		sessionRepo:     sessionRepo,
		processor:       processor,
		reportBuilder:   reportBuilder,
		reportFormatter: reportFormatter,
	}
}

func (e *Engine) Process(event events.Event) ([]events.Event, error) {
	if event.ID() == events.UserRegisteredEventID {
		_, _ = e.sessionRepo.Create(event.GetUserID())
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
		if err == models.ErrSessionNotFound {
			session, createErr := e.sessionRepo.Create(event.GetUserID())
			if createErr != nil {
				return nil, createErr
			}

			session.State = models.PlayerStateDisqualified
			if saveErr := e.sessionRepo.Save(session); saveErr != nil {
				return nil, saveErr
			}

			return []events.Event{
				events.UserDisqualifiedEvent{
					BaseEvent: events.BaseEvent{
						Timestamp: event.GetTimestamp(),
						UserID:    event.GetUserID(),
					},
				},
			}, nil
		}

		return nil, err
	}

	if event.ID() == events.UserLeftDungeonEventID {
		out, err := e.processor.Process(session, event)
		if err != nil {
			return nil, err
		}

		if err := e.sessionRepo.Save(session); err != nil {
			return nil, err
		}

		return out, nil
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
				EventID: event.GetEventID(),
			},
		}, nil
	default:
	}

	out, err := e.processor.Process(session, event)
	if err != nil {
		return nil, err
	}

	if err := e.sessionRepo.Save(session); err != nil {
		return nil, err
	}

	return out, nil
}

func (e *Engine) Report() string {
	if e.reportBuilder == nil || e.reportFormatter == nil {
		return "Final report:"
	}

	sessions, err := e.sessionRepo.GetAll()
	if err != nil {
		return "error generating report"
	}

	userReports := e.reportBuilder.Build(sessions)
	return e.reportFormatter.Format(userReports)
}
