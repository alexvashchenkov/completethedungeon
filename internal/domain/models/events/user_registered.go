package events

type UserRegisteredEvent struct {
	BaseEvent
}

func (UserRegisteredEvent) ID() EventID {
	return UserRegisteredEventID
}
