package events

type UserEnteredDungeonEvent struct {
	BaseEvent
}

func (UserEnteredDungeonEvent) ID() EventID {
	return UserEnteredDungeonEventID
}
