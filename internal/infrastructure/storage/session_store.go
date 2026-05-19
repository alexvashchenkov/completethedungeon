package storage

import (
	"impulse/internal/domain/models"
)

type InMemorySessionStore struct {
	sessions map[int]*models.PlayerSession
}

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		sessions: make(map[int]*models.PlayerSession),
	}
}

func (s *InMemorySessionStore) GetAll() ([]*models.PlayerSession, error) {
	out := make([]*models.PlayerSession, 0, len(s.sessions))
	for _, session := range s.sessions {
		out = append(out, session)
	}
	return out, nil
}

func (s *InMemorySessionStore) Get(userID int) (*models.PlayerSession, error) {
	session, exists := s.sessions[userID]
	if !exists {
		return nil, models.ErrSessionNotFound
	}
	return session, nil
}

func (s *InMemorySessionStore) Create(userID int) (*models.PlayerSession, error) {
	session, exists := s.sessions[userID]
	if exists {
		return session, nil
	}

	s.sessions[userID] = &models.PlayerSession{
		ID:    userID,
		HP:    100,
		State: models.PlayerStateRegistered,
	}
	return s.sessions[userID], nil
}

func (s *InMemorySessionStore) Save(session *models.PlayerSession) error {
	s.sessions[session.ID] = session
	return nil
}
