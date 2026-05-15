package events

type UserEnteredBossFloorEvent struct {
	BaseEvent
}

func (UserEnteredBossFloorEvent) ID() EventID {
	return UserEnteredBossFloorEventID
}
