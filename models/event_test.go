package models

import (
	"strconv"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		e        Event
		expected error
	}{
		{
			Event{
				Type: Birthday,
			},
			ErrMissingNameOrNames,
		},
		{
			Event{
				Name:  "John",
				Names: [2]string{"John"},
				Type:  Birthday,
			},
			ErrNameOrNames,
		},
		{
			Event{
				Name:  "John",
				Names: [2]string{"John", "Jane"},
				Type:  Birthday,
			},
			ErrNameOrNames,
		},
		{
			Event{
				Names: [2]string{"John"},
				Type:  Birthday,
			},
			ErrNamesArePair,
		},
		{
			Event{
				Names: [2]string{"John", "Jane"},
				Type:  "asdf",
			},
			ErrInvalidEventType,
		},
		{
			Event{
				Name: "John",
				Type: Birthday,
			},
			nil,
		},
		{
			Event{
				Names: [2]string{"John", "Jane"},
				Type:  Wedding,
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if err := tt.e.Validate(); err != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, err)
			}
		})
	}
}
