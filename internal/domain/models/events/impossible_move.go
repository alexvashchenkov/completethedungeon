package events

type UserMakesImpossibleMove struct {
	BaseEvent
	EventID int `json:"event_id"`
}

func (UserMakesImpossibleMove) ID() EventID {
	return UserMakesImpossibleMoveEventID
}
