package events

import "time"

type EventID int

const (
	BaseEventID EventID = iota
	UserRegisteredEventID
	UserEnteredDungeonEventID
	UserKilledMonsterEventID
	UserWentToNextFloorEventID
	UserWentToPreviousFloorEventID
	UserEnteredBossFloorEventID
	UserKilledBossEventID
	UserLeftDungeonEventID
	UserCantContinueDueToReasonEventID
	UserRestoredHealthEventID
	UserReceivedDamageEventID
)

const (
	UserDisqualifiedEventID        EventID = 31
	UserDiedEventID                EventID = 32
	UserMakesImpossibleMoveEventID EventID = 33
)

type Event interface {
	ID() EventID
	GetEventID() int
	GetTimestamp() time.Time
	GetUserID() int
}

type BaseEvent struct {
	EventID   int       `json:"event_id"`
	Timestamp time.Time `json:"timestamp"`
	UserID    int       `json:"user_id"`
}

func (BaseEvent) ID() EventID {
	return BaseEventID
}

func (b BaseEvent) GetTimestamp() time.Time {
	return b.Timestamp
}
func (b BaseEvent) GetUserID() int {
	return b.UserID
}
func (b BaseEvent) GetEventID() int {
	return b.EventID
}
