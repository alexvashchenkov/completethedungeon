package events

type UserLeftDungeonEvent struct {
	BaseEvent
}

func (UserLeftDungeonEvent) ID() EventID {
	return UserLeftDungeonEventID
}
