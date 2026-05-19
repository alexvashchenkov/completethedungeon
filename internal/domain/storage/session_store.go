package storage

import "impulse/internal/domain/models"

type SessionStore interface {
	GetAll() ([]*models.PlayerSession, error)
	Get(userID int) (*models.PlayerSession, error)
	Create(userID int) (*models.PlayerSession, error)
	Save(session *models.PlayerSession) error
}
