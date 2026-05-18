package parser

import (
	"reflect"
	"testing"
)

func TestTokenizer_Tokenize(t *testing.T) {
	tokens := NewTokenizer()

	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "without extra param",
			input: "[14:00:00] 1 1",
			want:  []string{"[14:00:00]", "1", "1"},
		},
		{
			name:  "with multiword extra param",
			input: "[14:00:00] 1 9 too tired to continue",
			want:  []string{"[14:00:00]", "1", "9", "too tired to continue"},
		},
		{
			name:    "invalid format",
			input:   "[14:00:00] 1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokens.Tokenize(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("unexpected tokens\nwant: %#v\ngot:  %#v", tt.want, got)
			}
		})
	}
}
