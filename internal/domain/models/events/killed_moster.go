package events

type UserKilledMonsterEvent struct {
	BaseEvent
}

func (UserKilledMonsterEvent) ID() EventID {
	return UserKilledMonsterEventID
}
