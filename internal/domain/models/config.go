package models

import "time"

type Config struct {
	Floors   int       `json:"Floors"`
	Monsters int       `json:"Monsters"`
	OpenAt   time.Time `json:"OpenAt"`
	Duration int       `json:"Duration"`
}
