package output

import (
	"fmt"
	"impulse/internal/adapters/output/formatter"
	"impulse/internal/domain/models/events"
	"io"
	"log"
)

type StdoutEventSink struct {
	writer    io.Writer
	formatter formatter.Formatter
}

func NewStdoutEventSink(w io.Writer) EventSink {
	return &StdoutEventSink{
		writer:    w,
		formatter: formatter.NewStringFormatter(),
	}
}

func (s StdoutEventSink) Write(event events.Event) error {
	_, err := fmt.Fprintln(
		s.writer,
		s.formatter.Format(event),
	)

	return err
}

func (s StdoutEventSink) WriteMany(events []events.Event) error {
	for _, event := range events {
		if err := s.Write(event); err != nil {
			log.Printf("unable to write event to output: %v\n", err)
		}
	}
	return nil
}
