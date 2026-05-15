package storage

import "impulse/internal/domain/models"

type SessionStore interface {
	Get(userID int) (*models.PlayerSession, error)
	Create(userID int) (*models.PlayerSession, error)
}
