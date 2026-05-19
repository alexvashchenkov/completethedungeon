package events

type UserDisqualifiedEvent struct {
	BaseEvent
}

func (UserDisqualifiedEvent) ID() EventID {
	return UserDisqualifiedEventID
}
