package config

import (
	"encoding/json"
	"impulse/internal/domain/models"
	"io"
)

type JsonConfigSource struct {
	reader io.Reader
}

func NewJsonConfigSource(r io.Reader) ConfigSource {
	return &JsonConfigSource{
		reader: r,
	}
}

func (j JsonConfigSource) Load() (models.Config, error) {
	decoder := json.NewDecoder(j.reader)

	var cfg models.Config

	err := decoder.Decode(&cfg)

	return cfg, err
}
