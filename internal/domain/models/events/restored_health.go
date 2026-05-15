package events

type UserRestoredHealthEvent struct {
	BaseEvent
	Amount int `json:"amount"`
}

func (UserRestoredHealthEvent) ID() EventID {
	return UserRestoredHealthEventID
}
