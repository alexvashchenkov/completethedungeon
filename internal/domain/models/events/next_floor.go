package events

type UserWentToNextFloorEvent struct {
	BaseEvent
}

func (UserWentToNextFloorEvent) ID() EventID {
	return UserWentToNextFloorEventID
}
