package report

import (
	"impulse/internal/domain/models"
	"sort"
	"time"
)

type ReportState int

const (
	SUCCESS ReportState = iota
	FAIL
	DISQUAL
)

func (s ReportState) String() string {
	switch s {
	case SUCCESS:
		return "SUCCESS"
	case FAIL:
		return "FAIL"
	case DISQUAL:
		return "DISQUAL"
	default:
		return "UNKNOWN"
	}
}

type FinalReport struct {
	State                ReportState
	UserID               int
	DungeonDuration      time.Duration
	AverageFloorDuration time.Duration
	BossDuration         time.Duration
	HP                   int
}

type Builder interface {
	Build(sessions []*models.PlayerSession) []*FinalReport
}

type ReportBuilder struct {
	config *models.Config
}

func NewReportBuilder(config *models.Config) Builder {
	return &ReportBuilder{
		config: config,
	}
}

func (b *ReportBuilder) Build(sessions []*models.PlayerSession) []*FinalReport {
	sort.Slice(sessions, func(i, j int) bool { return sessions[i].ID < sessions[j].ID })

	out := make([]*FinalReport, 0)

	for _, session := range sessions {
		out = append(out, buildSingleReport(session, b.config))
	}

	return out
}

func buildSingleReport(session *models.PlayerSession, config *models.Config) *FinalReport {
	dungeonDuration := session.Metrics.TotalDungeonDuration
	if dungeonDuration == 0 && !session.DungeonEnteredAt.IsZero() && !session.DungeonFinishedAt.IsZero() {
		dungeonDuration = session.DungeonFinishedAt.Sub(session.DungeonEnteredAt)
	}

	out := &FinalReport{
		State:                determineReportState(session, config),
		UserID:               session.ID,
		HP:                   session.HP,
		DungeonDuration:      dungeonDuration,
		AverageFloorDuration: calculateAvgFloorTime(session.Metrics.FloorDurations),
		BossDuration:         session.Metrics.BossDuration,
	}

	return out
}

func determineReportState(session *models.PlayerSession, config *models.Config) ReportState {
	if session.State == models.PlayerStateDisqualified {
		return DISQUAL
	} else if session.BossKilled && (session.State == models.PlayerStateCompleted || session.State == models.PlayerStateLeftDungeon) {
		return SUCCESS
	}
	return FAIL
}

func calculateAvgFloorTime(floorDurations []time.Duration) time.Duration {
	if len(floorDurations) == 0 {
		return 0
	}

	sum := float64(0)

	for _, d := range floorDurations {
		sum += d.Seconds()
	}

	avg := sum / float64(len(floorDurations))
	return time.Duration(avg * float64(time.Second)).Round(time.Second)
}
