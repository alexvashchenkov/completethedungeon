package parser

import (
	"fmt"
	"impulse/internal/domain/models/events"
	"strconv"
)

type EventFactory struct{}

func NewEventFactory() *EventFactory {
	return &EventFactory{}
}

func (e *EventFactory) Build(rawEvent *RawEvent) (events.Event, error) {
	switch events.EventID(rawEvent.EventID) {
	case events.UserRegisteredEventID:
		return events.UserRegisteredEvent{
			BaseEvent: events.BaseEvent{
				Timestamp: rawEvent.Timestamp,
				UserID:    rawEvent.UserID,
			},
		}, nil
	case events.UserEnteredDungeonEventID:
		return events.UserEnteredDungeonEvent{
			BaseEvent: events.BaseEvent{
				Timestamp: rawEvent.Timestamp,
				UserID:    rawEvent.UserID,
			},
		}, nil

	case events.UserKilledMonsterEventID:
		return events.UserKilledMonsterEvent{
			BaseEvent: events.BaseEvent{
				Timestamp: rawEvent.Timestamp,
				UserID:    rawEvent.UserID,
			},
		}, nil

	case events.UserWentToNextFloorEventID:
		return events.UserWentToNextFloorEvent{
			BaseEvent: events.BaseEvent{
				Timestamp: rawEvent.Timestamp,
				UserID:    rawEvent.UserID,
			},
		}, nil

	case events.UserWentToPreviousFloorEventID:
		return events.UserWentToPreviousFloorEvent{
			BaseEvent: events.BaseEvent{
				Timestamp: rawEvent.Timestamp,
				UserID:    rawEvent.UserID,
			},
		}, nil

	case events.UserEnteredBossFloorEventID:
		return events.UserEnteredBossFloorEvent{
			BaseEvent: events.BaseEvent{
				Timestamp: rawEvent.Timestamp,
				UserID:    rawEvent.UserID,
			},
		}, nil

	case events.UserKilledBossEventID:
		return events.UserKilledBossEvent{
			BaseEvent: events.BaseEvent{
				Timestamp: rawEvent.Timestamp,
				UserID:    rawEvent.UserID,
			},
		}, nil

	case events.UserLeftDungeonEventID:
		return events.UserLeftDungeonEvent{
			BaseEvent: events.BaseEvent{
				Timestamp: rawEvent.Timestamp,
				UserID:    rawEvent.UserID,
			},
		}, nil

	case events.UserCantContinueDueToReasonEventID:
		return events.UserCantContinueDueToReasonEvent{
			BaseEvent: events.BaseEvent{
				Timestamp: rawEvent.Timestamp,
				UserID:    rawEvent.UserID,
			},
			Reason: rawEvent.Extra,
		}, nil

	case events.UserRestoredHealthEventID:
		amount, err := strconv.Atoi(rawEvent.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to parse health amount: %w", err)
		}
		return events.UserRestoredHealthEvent{
			BaseEvent: events.BaseEvent{
				Timestamp: rawEvent.Timestamp,
				UserID:    rawEvent.UserID,
			},
			Amount: amount,
		}, nil

	case events.UserReceivedDamageEventID:
		amount, err := strconv.Atoi(rawEvent.Extra)
		if err != nil {
			return nil, fmt.Errorf("failed to parse damage amount: %w", err)
		}

		return events.UserReceivedDamageEvent{
			BaseEvent: events.BaseEvent{
				Timestamp: rawEvent.Timestamp,
				UserID:    rawEvent.UserID,
			},
			Amount: amount,
		}, nil
	}
	return nil, fmt.Errorf("unknown event: %v", events.EventID(rawEvent.EventID))
}
