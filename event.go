package main

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

const (
	birthday eventType = "birthday"
	nameday  eventType = "nameday"
	wedding  eventType = "wedding anniversary"
)

type eventType string

var _ fmt.Stringer = eventType("")

type event struct {
	Name  *string    `yaml:"name"`
	Names *[2]string `yaml:"names"`
	Type  eventType  `yaml:"type"`
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

func (et *eventType) UnmarshalYAML(value *yaml.Node) error {
	switch v := eventType(value.Value); v {
	case birthday, nameday, wedding:
		*et = v
	default:
		return fmt.Errorf("invalid type: %s", value.Value)
	}
	return nil
}

func (e event) String() string {
	if name := e.Name; name != nil {
		return fmt.Sprintf("%s ma dziś %s!", *name, e.Type)
	}
	return fmt.Sprintf("%s i %s mają dziś %s!", e.Names[0], e.Names[1], e.Type)
}
