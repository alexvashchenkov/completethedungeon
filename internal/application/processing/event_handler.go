package processing

import (
	"impulse/internal/domain/models"
	"impulse/internal/domain/models/events"
)

type EventHandler interface {
	Handle(session *models.PlayerSession, event events.Event) ([]events.Event, error)
}
