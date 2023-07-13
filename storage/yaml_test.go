package storage

import (
	"strconv"
	"testing"

	"git.sr.ht/~tymek/przypominajka/models"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		e        models.Event
		expected error
	}{
		{
			models.Event{
				Type: models.Birthday,
			},
			errMissingNameOrNames,
		},
		{
			models.Event{
				Name:  "John",
				Names: [2]string{"John"},
				Type:  models.Birthday,
			},
			errNameOrNames,
		},
		{
			models.Event{
				Name:  "John",
				Names: [2]string{"John", "Jane"},
				Type:  models.Birthday,
			},
			errNameOrNames,
		},
		{
			models.Event{
				Names: [2]string{"John"},
				Type:  models.Birthday,
			},
			errNamesArePair,
		},
		{
			models.Event{
				Names: [2]string{"John", "Jane"},
				Type:  "asdf",
			},
			errInvalidEventType,
		},
		{
			models.Event{
				Name: "John",
				Type: models.Birthday,
			},
			nil,
		},
		{
			models.Event{
				Names: [2]string{"John", "Jane"},
				Type:  models.Wedding,
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if err := validate(tt.e); err != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, err)
			}
		})
	}
}
