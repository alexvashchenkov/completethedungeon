package models

import "time"

type PlayerState int

const (
	PlayerStateRegistered   PlayerState = 1
	PlayerStateInDungeon    PlayerState = 2
	PlayerStateDisqualified PlayerState = 3
	PlayerStateDead         PlayerState = 4
	PlayerStateLeftDungeon  PlayerState = 5
	PlayerStateCompleted    PlayerState = 6
)

type PlayerSession struct {
	ID int

	State PlayerState

	HP int

	CurrentFloor int

	CurrentFloorStartedAt time.Time

	CurrentFloorFinishedAt time.Time

	MonstersKilledOnFloor int

	DungeonEnteredAt time.Time

	DungeonFinishedAt time.Time

	Metrics Metrics
}

type Metrics struct {
	FloorDurations []time.Duration

	BossDuration time.Duration

	TotalDungeonDuration time.Duration
}
