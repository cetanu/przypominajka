package main

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
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

func newErrInvalidEventType(et eventType) error {
	return fmt.Errorf("%w: %s", errInvalidEventType, et)
}

type eventType string

var _ fmt.Stringer = eventType("")

type event struct {
	Name    string    `yaml:"name"`
	Names   [2]string `yaml:"names"`
	Surname string    `yaml:"surname"`
	Type    eventType `yaml:"type"`
}

var (
	_ fmt.Stringer     = event{}
	_ yaml.Unmarshaler = (*event)(nil)
)

func (et eventType) String() string {
	switch et {
	case birthday:
		return formatBirthday
	case nameday:
		return formatNameday
	case wedding:
		return formatWeddingAnniversary
	}
	return string(et)
}

func (e event) String() string {
	if e.Name != "" {
		if e.Surname != "" {
			return fmt.Sprintf(formatSingularSurname, e.Name, e.Surname, e.Type)
		}
		return fmt.Sprintf(formatSingular, e.Name, e.Type)
	}
	if e.Surname != "" {
		return fmt.Sprintf(formatMessagePluralSurname, e.Names[0], e.Names[1], e.Surname, e.Type)
	}
	return fmt.Sprintf(formatMessagePlural, e.Names[0], e.Names[1], e.Type)
}

func (e *event) UnmarshalYAML(value *yaml.Node) error {
	s := struct {
		Name    string    `yaml:"name"`
		Names   [2]string `yaml:"names"`
		Surname string    `yaml:"surname"`
		Type    eventType `yaml:"type"`
	}{}
	if err := value.Decode(&s); err != nil {
		return err
	}
	*e = s
	return e.validate()
}

func (e event) format(list bool) string {
	switch {
	case e.Name != "" && e.Surname == "" && !list:
		return fmt.Sprintf(formatSingular, e.Name, e.Type)
	case e.Name != "" && e.Surname == "":
		return fmt.Sprintf(formatListSingular, e.Name, e.Type)
	case e.Name != "" && e.Surname != "" && !list:
		return fmt.Sprintf(formatSingularSurname, e.Name, e.Surname, e.Type)
	case e.Name != "" && e.Surname != "":
		return fmt.Sprintf(formatListSingularSurname, e.Name, e.Surname, e.Type)
		// Plural
	case e.Surname == "" && !list:
		return fmt.Sprintf(formatMessagePlural, e.Names[0], e.Names[1], e.Type)
	case e.Surname == "" && list:
		return fmt.Sprintf(formatListMessagePlural, e.Names[0], e.Names[1], e.Type)
	case e.Surname != "" && !list:
		return fmt.Sprintf(formatMessagePluralSurname, e.Names[0], e.Names[1], e.Surname, e.Type)
	case e.Surname != "" && list:
		return fmt.Sprintf(formatListMessagePluralSurname, e.Names[0], e.Names[1], e.Surname, e.Type)
	}
	return ""
}

func (e event) validate() error {
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
		return newErrInvalidEventType(e.Type)
	}
	return nil
}
