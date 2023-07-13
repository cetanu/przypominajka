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
				Type: birthday,
			},
			errMissingNameOrNames,
		},
		{
			Event{
				Name:  "John",
				Names: [2]string{"John"},
				Type:  birthday,
			},
			errNameOrNames,
		},
		{
			Event{
				Name:  "John",
				Names: [2]string{"John", "Jane"},
				Type:  birthday,
			},
			errNameOrNames,
		},
		{
			Event{
				Names: [2]string{"John"},
				Type:  birthday,
			},
			errNamesArePair,
		},
		{
			Event{
				Names: [2]string{"John", "Jane"},
				Type:  "asdf",
			},
			errInvalidEventType,
		},
		{
			Event{
				Name: "John",
				Type: birthday,
			},
			nil,
		},
		{
			Event{
				Names: [2]string{"John", "Jane"},
				Type:  wedding,
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if err := tt.e.Validate(); err != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, err)
			}
		})
	}
}
