package formatter

import (
	"fmt"
	"impulse/internal/domain/models/events"
)

type StringFormatter struct{}

func NewStringFormatter() Formatter {
	return &StringFormatter{}
}

func (sf *StringFormatter) Format(ev events.Event) string {
	var info string
	timestamp := ev.GetTimestamp().Format("15:04:05")

	switch v := ev.(type) {

	case events.UserRegisteredEvent:
		info = "registered"
	case events.UserEnteredDungeonEvent:
		info = "entered the dungeon"
	case events.UserKilledMonsterEvent:
		info = "killed the monster"
	case events.UserWentToNextFloorEvent:
		info = "went to the next floor"
	case events.UserWentToPreviousFloorEvent:
		info = "went to the previous floor"
	case events.UserEnteredBossFloorEvent:
		info = "entered the boss's floor"
	case events.UserKilledBossEvent:
		info = "killed the boss"
	case events.UserLeftDungeonEvent:
		info = "left the dungeon"
	case events.UserCantContinueDueToReasonEvent:
		info = fmt.Sprintf("cant continue due to [%s]", v.Reason)
	case events.UserRestoredHealthEvent:
		info = fmt.Sprintf("has restored [%d] of health", v.Amount)
	case events.UserReceivedDamageEvent:
		info = fmt.Sprintf("recieved [%d] of damage", v.Amount)
	case events.UserDiedEvent:
		info = "is dead"
	case events.UserDisqualifiedEvent:
		info = "is disqualified"
	case events.UserMakesImpossibleMove:
		info = fmt.Sprintf("makes impossible move [%d]", v.EventID)
	}

	return fmt.Sprintf("[%s] Player [%d] %s", timestamp, ev.GetUserID(), info)
}
