package events

type UserMakesImpossibleMove struct {
	BaseEvent
	EventID EventID `json:"event_id"`
}

func (UserMakesImpossibleMove) ID() EventID {
	return UserMakesImpossibleMoveEventID
}
