package main

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type info struct {
	Name  *string    `yaml:"name"`
	Names *[2]string `yaml:"names"`
	Type  eventType  `yaml:"type"`
}

var _ fmt.Stringer = info{}

func (i info) String() string {
	if name := i.Name; name != nil {
		return fmt.Sprintf("%s ma dziś %s", *name, i.Type)
	}
	return fmt.Sprintf("%s i %s mają dziś %s", i.Names[0], i.Names[1], i.Type)
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
		return "imieniny"
	case wedding:
		return "rocznicę ślubu"
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
