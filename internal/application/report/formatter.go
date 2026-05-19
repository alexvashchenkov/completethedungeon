package report

import (
	"fmt"
	"strings"
	"time"
)

type Formatter interface {
	Format(reports []*FinalReport) string
	FormatOne(report *FinalReport) string
}

type StringFormatter struct {
}

func NewStringFormatter() Formatter {
	return &StringFormatter{}
}

func (f *StringFormatter) Format(reports []*FinalReport) string {
	output := "Final report:\n"
	for _, report := range reports {
		output += f.FormatOne(report) + "\n"
	}

	return strings.TrimSuffix(output, "\n")
}

func (f *StringFormatter) FormatOne(report *FinalReport) string {
	return fmt.Sprintf("[%s] %d [%s, %s, %s] HP:%d",
		report.State,
		report.UserID,
		formatDuration(report.DungeonDuration),
		formatDuration(report.AverageFloorDuration),
		formatDuration(report.BossDuration),
		report.HP,
	)
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}

	totalSeconds := int64(d / time.Second)
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
