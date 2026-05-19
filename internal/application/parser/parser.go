package parser

import (
	"fmt"
	"impulse/internal/domain/models/events"
	"strconv"
	"strings"
	"time"
)

type Parser struct {
	tokenizer *Tokenizer
	rawParser *RawParser
	factory   *EventFactory
}

func NewParser() *Parser {
	return &Parser{
		tokenizer: NewTokenizer(),
		rawParser: NewRawParser(),
		factory:   NewEventFactory(),
	}
}

func (p *Parser) Parse(line string) (events.Event, error) {
	tokens, err := p.tokenizer.Tokenize(line)
	if err != nil {
		return nil, fmt.Errorf("unable to tokenize: %w", err)
	}

	raw, err := p.rawParser.Parse(tokens)
	if err != nil {
		return nil, fmt.Errorf("unable to parse raw event: %w", err)
	}

	return p.factory.Build(raw)
}

type RawParser struct{}

func NewRawParser() *RawParser {
	return &RawParser{}
}

func (p *RawParser) Parse(tokens []string) (*RawEvent, error) {
	timestamp, err := time.Parse("15:04:05", strings.Trim(tokens[0], "[]"))
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp: %w", err)
	}

	userID, err := strconv.Atoi(tokens[1])
	if err != nil {
		return nil, fmt.Errorf("invalid event id: %w", err)
	}

	eventID, err := strconv.Atoi(tokens[2])
	if err != nil {
		return nil, fmt.Errorf("invalid event id: %w", err)
	}

	output := &RawEvent{
		Timestamp: timestamp,
		EventID:   eventID,
		UserID:    userID,
	}

	if len(tokens) > 3 {
		output.Extra = strings.Join(tokens[3:], " ")
	}

	return output, nil
}
