package main

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type info struct {
	Name  string    `yaml:"name"`
	Names []string  `yaml:"names"`
	Type  eventType `yaml:"type"`
}

const (
	birthday eventType = "birthday"
	nameday  eventType = "nameday"
	wedding  eventType = "wedding anniversary"
)

type eventType string

var _ fmt.Stringer = eventType("")

func (e eventType) String() string {
	switch e {
	case birthday:
		return "urodziny"
	case nameday:
		return "nameday"
	default:
		return "unknown"
	}
}

func (e *eventType) UnmarshalYAML(value *yaml.Node) error {
	switch v := eventType(value.Value); v {
	case birthday, nameday, wedding:
		*e = v
	default:
		return fmt.Errorf("invalid type: %s", value.Value)
	}
	return nil
}
