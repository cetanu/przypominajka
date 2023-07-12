package main

import (
	"strconv"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		e        event
		expected error
	}{
		{
			event{
				Type: birthday,
			},
			errMissingNameOrNames,
		},
		{
			event{
				Name:  "John",
				Names: [2]string{"John"},
				Type:  birthday,
			},
			errNameOrNames,
		},
		{
			event{
				Name:  "John",
				Names: [2]string{"John", "Jane"},
				Type:  birthday,
			},
			errNameOrNames,
		},
		{
			event{
				Names: [2]string{"John"},
				Type:  birthday,
			},
			errNamesArePair,
		},
		{
			event{
				Names: [2]string{"John", "Jane"},
				Type:  "asdf",
			},
			errInvalidEventType,
		},
		{
			event{
				Name: "John",
				Type: birthday,
			},
			nil,
		},
		{
			event{
				Names: [2]string{"John", "Jane"},
				Type:  wedding,
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if err := tt.e.validate(); err != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, err)
			}
		})
	}
}
