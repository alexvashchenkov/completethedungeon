package events

type UserWentToPreviousFloorEvent struct {
	BaseEvent
}

func (UserWentToPreviousFloorEvent) ID() EventID {
	return UserWentToPreviousFloorEventID
}
