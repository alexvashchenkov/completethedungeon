package events

type UserReceivedDamageEvent struct {
	BaseEvent
	Amount int `json:"amount"`
}

func (UserReceivedDamageEvent) ID() EventID {
	return UserReceivedDamageEventID
}
