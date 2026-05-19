package config

import "impulse/internal/domain/models"

type ConfigSource interface {
	Load() (models.Config, error)
}
