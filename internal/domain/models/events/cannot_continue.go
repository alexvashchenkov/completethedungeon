package events

type UserCantContinueDueToReasonEvent struct {
	BaseEvent
	Reason string `json:"reason"`
}

func (UserCantContinueDueToReasonEvent) ID() EventID {
	return UserCantContinueDueToReasonEventID
}
