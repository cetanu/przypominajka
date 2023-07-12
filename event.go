package main

import (
	"errors"
	"fmt"
)

const (
	birthday eventType = "birthday"
	nameday  eventType = "nameday"
	wedding  eventType = "wedding anniversary"
)

var (
	errMissingNameOrNames = errors.New("'name' or 'names' must be provided")
	errNameOrNames        = errors.New("'name' is mutually exclusive with 'names'")
	errNamesArePair       = errors.New("'names' must have two elements")
	errInvalidEventType   = errors.New("invalid event type")
)

type eventType string

var _ fmt.Stringer = eventType("")

type event struct {
	Name    string    `yaml:"name"`
	Names   [2]string `yaml:"names"`
	Surname string    `yaml:"surname"`
	Type    eventType `yaml:"type"`
}

var _ fmt.Stringer = event{}

func (et eventType) String() string {
	switch et {
	case birthday:
		return "urodziny"
	case nameday:
		return "imieniny"
	case wedding:
		return "rocznicę ślubu"
	default:
		return "unknown"
	}
}

func (e event) String() string {
	if e.Name != "" {
		if e.Surname != "" {
			return fmt.Sprintf("%s %s ma dziś %s!", e.Name, e.Surname, e.Type)
		}
		return fmt.Sprintf("%s ma dziś %s!", e.Name, e.Type)
	}
	if e.Surname != "" {
		return fmt.Sprintf("%s i %s %s mają dziś %s!", e.Names[0], e.Names[1], e.Surname, e.Type)
	}
	return fmt.Sprintf("%s i %s mają dziś %s!", e.Names[0], e.Names[1], e.Type)
}

func (e event) Validate() error {
	if e.Name == "" && (e.Names[0] == "" && e.Names[1] == "") {
		return errMissingNameOrNames
	}
	if e.Name != "" && (e.Names[0] != "" || e.Names[1] != "") {
		return errNameOrNames
	}
	if e.Name == "" && (e.Names[0] == "" || e.Names[1] == "") {
		return errNamesArePair
	}
	switch e.Type {
	case birthday, nameday, wedding:
	default:
		return errInvalidEventType
	}

	return nil
}
