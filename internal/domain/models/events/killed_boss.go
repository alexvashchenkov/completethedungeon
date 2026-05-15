package events

type UserKilledBossEvent struct {
	BaseEvent
}

func (UserKilledBossEvent) ID() EventID {
	return UserKilledBossEventID
}
