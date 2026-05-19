package events

import (
	"bufio"
	"impulse/internal/application/parser"
	"impulse/internal/domain/models/events"
	"io"
)

type FileEventSource struct {
	scanner *bufio.Scanner
	parser  *parser.Parser
}

func NewFileEventSource(r io.Reader) EventSource {

	return &FileEventSource{
		scanner: bufio.NewScanner(r),
		parser:  parser.NewParser(),
	}
}

func (s *FileEventSource) Next() (events.Event, error) {
	if !s.scanner.Scan() {
		if err := s.scanner.Err(); err != nil {
			return nil, err
		}
		return nil, io.EOF
	}
	line := s.scanner.Text()

	return s.parser.Parse(line)
}
