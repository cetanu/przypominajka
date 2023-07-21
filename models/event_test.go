package models

import (
	"strconv"
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		e        Event
		list     bool
		expected string
	}{
		{
			Event{
				Name: "John",
				Type: Birthday,
			},
			false,
			"John obchodzi dziś urodziny!",
		},
		{
			Event{
				Name:    "John",
				Surname: "Doe",
				Type:    Birthday,
			},
			false,
			"John Doe obchodzi dziś urodziny!",
		},
		{
			Event{
				Name:    "John",
				Surname: "Doe",
				Type:    Birthday,
				Month:   time.February,
				Day:     15,
			},
			true,
			"15.02 - John Doe obchodzi urodziny",
		},
		{
			Event{
				Name:    "John",
				Surname: "Doe",
				Type:    Birthday,
				Month:   time.February,
				Day:     15,
			},
			true,
			"15.02 - John Doe obchodzi urodziny",
		},
		{
			Event{
				Names:   &[2]string{"Jane", "John"},
				Surname: "Doe",
				Type:    Wedding,
			},
			false,
			"Jane i John Doe obchodzą dziś rocznicę ślubu!",
		},
		{
			Event{
				Names:   &[2]string{"Jane", "John"},
				Surname: "Doe",
				Type:    Wedding,
				Month:   time.February,
				Day:     15,
			},
			true,
			"15.02 - Jane i John Doe obchodzą rocznicę ślubu",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if err := tt.e.Format(tt.list); err != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, err)
			}
		})
	}
}

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
				Names: &[2]string{"John"},
				Type:  Birthday,
			},
			ErrNameOrNames,
		},
		{
			Event{
				Name:  "John",
				Names: &[2]string{"John", "Jane"},
				Type:  Birthday,
			},
			ErrNameOrNames,
		},
		{
			Event{
				Names: &[2]string{"John"},
				Type:  Birthday,
			},
			ErrNamesArePair,
		},
		{
			Event{
				Names: &[2]string{"John", "Jane"},
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
				Names: &[2]string{"John", "Jane"},
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
