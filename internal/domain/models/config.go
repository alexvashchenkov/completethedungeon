package models

type Config struct {
	Floors   int `json:"Floors"`
	Monsters int `json:"Monsters"`
	OpenAt   int `json:"OpenAt"`
	Duration int `json:"Duration"`
}
