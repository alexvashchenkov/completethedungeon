package parser

import (
	"errors"
	"strings"
)

var (
	ErrInvalidFormat = errors.New("invalid format")
)

type Tokenizer struct{}

func NewTokenizer() *Tokenizer {
	return &Tokenizer{}
}

func (t *Tokenizer) Tokenize(text string) ([]string, error) {
	parts := strings.Fields(text)

	if len(parts) < 3 {
		return nil, ErrInvalidFormat
	} else if len(parts) == 3 {
		return parts, nil
	}

	output := make([]string, 0, len(parts))
	output = append(output, parts[0])
	output = append(output, parts[1])
	output = append(output, parts[2])

	if len(parts) > 3 {
		output = append(output, strings.Join(parts[3:], " "))
	}

	return output, nil
}
