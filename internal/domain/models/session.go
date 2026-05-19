package models

import "time"

type PlayerState int

const (
	PlayerStateRegistered PlayerState = iota
	PlayerStateInDungeon
	PlayerStateDisqualified
	PlayerStateDead
	PlayerStateLeftDungeon
	PlayerStateCompleted
)

type PlayerSession struct {
	ID int

	State PlayerState

	HP int

	DungeonEnteredAt time.Time

	CurrentFloor int

	CurrentFloorStartedAt time.Time

	CurrentFloorFinishedAt time.Time

	MonstersKilledOnFloor int

	BossKilled bool

	BossKilledAt time.Time

	DungeonFinishedAt time.Time

	DungeonLeftAt time.Time

	Metrics Metrics
}

type Metrics struct {
	FloorsCompleted int

	FloorDurations []time.Duration

	BossDuration time.Duration

	TotalDungeonDuration time.Duration
}
