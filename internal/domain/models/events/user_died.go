package events

type UserDiedEvent struct {
	BaseEvent
}

func (UserDiedEvent) ID() EventID {
	return UserDiedEventID
}
