package events

import "time"

type EventID int

const (
	BaseEventID                        EventID = 0
	UserRegisteredEventID              EventID = 1
	UserEnteredDungeonEventID          EventID = 2
	UserKilledMonsterEventID           EventID = 3
	UserWentToNextFloorEventID         EventID = 4
	UserWentToPreviousFloorEventID     EventID = 5
	UserEnteredBossFloorEventID        EventID = 6
	UserKilledBossEventID              EventID = 7
	UserLeftDungeonEventID             EventID = 8
	UserCantContinueDueToReasonEventID EventID = 9
	UserRestoredHealthEventID          EventID = 10
	UserReceivedDamageEventID          EventID = 11
)

const (
	UserDisqualifiedEventID        EventID = 31
	UserDiedEventID                EventID = 32
	UserMakesImpossibleMoveEventID EventID = 33
)

type Event interface {
	ID() EventID
	GetTimestamp() time.Time
	GetUserID() int
}

type BaseEvent struct {
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
