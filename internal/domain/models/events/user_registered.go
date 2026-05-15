package events

type UserRegisteredUserEvent struct {
	BaseEvent
}

func (UserRegisteredUserEvent) ID() EventID {
	return UserRegisteredEventID
}
