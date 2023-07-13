package models

import (
	"fmt"
	"strings"
	"time"
)

type Events []Event

var _ fmt.Stringer = Events(nil)

func (ev Events) String() string {
	if len(ev) == 0 {
		return formatNoEvents
	}
	lines := make([]string, len(ev))
	for i, e := range ev {
		lines[i] = e.Format(true)
	}
	return strings.Join(lines, "\n")
}

type Event struct {
	Name    string     `yaml:"name"`
	Names   [2]string  `yaml:"names"`
	Surname string     `yaml:"surname"`
	Type    EventType  `yaml:"type"`
	Month   time.Month `yaml:"-"`
	Day     int        `yaml:"-"`
}

func (e Event) Format(list bool) string {
	var result string
	switch {
	case e.Name != "" && e.Surname == "" && !list:
		result = fmt.Sprintf(formatSingular, e.Name, e.Type)
	case e.Name != "" && e.Surname == "" && list:
		result = fmt.Sprintf(formatListSingular, e.Name, e.Type)
	case e.Name != "" && e.Surname != "" && !list:
		result = fmt.Sprintf(formatSingularSurname, e.Name, e.Surname, e.Type)
	case e.Name != "" && e.Surname != "" && list:
		result = fmt.Sprintf(formatListSingularSurname, e.Name, e.Surname, e.Type)
	// Plural
	case e.Surname == "" && !list:
		result = fmt.Sprintf(formatMessagePlural, e.Names[0], e.Names[1], e.Type)
	case e.Surname == "" && list:
		result = fmt.Sprintf(formatListMessagePlural, e.Names[0], e.Names[1], e.Type)
	case e.Surname != "" && !list:
		result = fmt.Sprintf(formatMessagePluralSurname, e.Names[0], e.Names[1], e.Surname, e.Type)
	case e.Surname != "" && list:
		result = fmt.Sprintf(formatListMessagePluralSurname, e.Names[0], e.Names[1], e.Surname, e.Type)
	}
	if list {
		result = fmt.Sprintf(formatListLine, e.Day, e.Month, result)
	}
	return result
}

func (e Event) Validate() error {
	if e.Name == "" && (e.Names[0] == "" && e.Names[1] == "") {
		return ErrMissingNameOrNames
	}
	if e.Name != "" && (e.Names[0] != "" || e.Names[1] != "") {
		return ErrNameOrNames
	}
	if e.Name == "" && (e.Names[0] == "" || e.Names[1] == "") {
		return ErrNamesArePair
	}
	switch e.Type {
	case Birthday, Nameday, Wedding:
	default:
		return ErrInvalidEventType
	}
	return nil
}
