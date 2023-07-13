package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"git.sr.ht/~tymek/przypominajka/format"
	"gopkg.in/yaml.v3"
)

var (
	errMissingNameOrNames = errors.New("'name' or 'names' must be provided")
	errNameOrNames        = errors.New("'name' is mutually exclusive with 'names'")
	errNamesArePair       = errors.New("'names' must have two elements")
	errInvalidEventType   = errors.New("invalid event type")
)

type Events []Event

func (ev Events) Format(month time.Month, day int) string {
	lines := make([]string, len(ev))
	for i, e := range ev {
		lines[i] = fmt.Sprintf(format.ListLine, day, month, e.Format(true))
	}
	return strings.Join(lines, "\n")
}

type Event struct {
	Name    string     `yaml:"name"`
	Names   [2]string  `yaml:"names"`
	Surname string     `yaml:"surname"`
	Type    eventType  `yaml:"type"`
	Month   time.Month `yaml:"-"`
	Day     int        `yaml:"-"`
}

var _ yaml.Unmarshaler = (*Event)(nil)

func (e *Event) UnmarshalYAML(value *yaml.Node) error {
	s := struct {
		Name    string     `yaml:"name"`
		Names   [2]string  `yaml:"names"`
		Surname string     `yaml:"surname"`
		Type    eventType  `yaml:"type"`
		Month   time.Month `yaml:"-"`
		Day     int        `yaml:"-"`
	}{}
	if err := value.Decode(&s); err != nil {
		return err
	}
	*e = s
	return e.Validate()
}

func (e Event) Format(list bool) string {
	switch {
	case e.Name != "" && e.Surname == "" && !list:
		return fmt.Sprintf(format.Singular, e.Name, e.Type)
	case e.Name != "" && e.Surname == "":
		return fmt.Sprintf(format.ListSingular, e.Name, e.Type)
	case e.Name != "" && e.Surname != "" && !list:
		return fmt.Sprintf(format.SingularSurname, e.Name, e.Surname, e.Type)
	case e.Name != "" && e.Surname != "":
		return fmt.Sprintf(format.ListSingularSurname, e.Name, e.Surname, e.Type)
		// Plural
	case e.Surname == "" && !list:
		return fmt.Sprintf(format.MessagePlural, e.Names[0], e.Names[1], e.Type)
	case e.Surname == "" && list:
		return fmt.Sprintf(format.ListMessagePlural, e.Names[0], e.Names[1], e.Type)
	case e.Surname != "" && !list:
		return fmt.Sprintf(format.MessagePluralSurname, e.Names[0], e.Names[1], e.Surname, e.Type)
	case e.Surname != "" && list:
		return fmt.Sprintf(format.ListMessagePluralSurname, e.Names[0], e.Names[1], e.Surname, e.Type)
	}
	return ""
}

func (e Event) Validate() error {
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
