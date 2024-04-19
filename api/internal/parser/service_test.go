package parser

import (
	"reflect"
	"testing"
)

func TestService_tokenize(t *testing.T) {
	s := NewService()

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "Simple Arithmetic",
			input: "5+5*3-2",
			want:  []string{"5", "+", "5", "*", "3", "-", "2"},
		},
		{
			name:  "With Decimal Points",
			input: "3.1415*10+0.0095",
			want:  []string{"3.1415", "*", "10", "+", "0.0095"},
		},
		{
			name:  "Negative Numbers",
			input: "-15+3*2",
			want:  []string{"-15", "+", "3", "*", "2"},
		},
		{
			name:  "No Spaces",
			input: "2*2",
			want:  []string{"2", "*", "2"},
		},
		{
			name:  "With Spaces",
			input: " 2 * 2 ",
			want:  []string{"2", "*", "2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.tokenize(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}
