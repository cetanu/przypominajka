package models

import (
	"fmt"
	"strings"
	"time"

	"git.sr.ht/~tymek/przypominajka/format"
)

type Events []Event

var _ fmt.Stringer = Events(nil)

func (ev Events) String() string {
	if len(ev) == 0 {
		return format.NoEvents
	}
	lines := make([]string, len(ev))
	for i, e := range ev {
		lines[i] = e.Format(true)
	}
	return strings.Join(lines, "\n")
}

type Event struct {
	Name    string     `yaml:"name,omitempty"`
	Names   *[2]string `yaml:"names,omitempty"`
	Surname string     `yaml:"surname,omitempty"`
	Type    EventType  `yaml:"type"`
	Month   time.Month `yaml:"-"`
	Day     int        `yaml:"-"`
}

func (e Event) Format(list bool) string {
	var result string
	switch {
	case e.Name != "" && e.Surname == "" && !list:
		result = fmt.Sprintf(format.Singular, e.Name, e.Type)
	case e.Name != "" && e.Surname == "" && list:
		result = fmt.Sprintf(format.ListSingular, e.Name, e.Type)
	case e.Name != "" && e.Surname != "" && !list:
		result = fmt.Sprintf(format.SingularSurname, e.Name, e.Surname, e.Type)
	case e.Name != "" && e.Surname != "" && list:
		result = fmt.Sprintf(format.ListSingularSurname, e.Name, e.Surname, e.Type)
	// Plural
	case e.Surname == "" && !list:
		result = fmt.Sprintf(format.Plural, e.Names[0], e.Names[1], e.Type)
	case e.Surname == "" && list:
		result = fmt.Sprintf(format.ListPlural, e.Names[0], e.Names[1], e.Type)
	case e.Surname != "" && !list:
		result = fmt.Sprintf(format.PluralSurname, e.Names[0], e.Names[1], e.Surname, e.Type)
	case e.Surname != "" && list:
		result = fmt.Sprintf(format.ListPluralSurname, e.Names[0], e.Names[1], e.Surname, e.Type)
	}
	if list {
		result = fmt.Sprintf(format.ListLine, e.Day, e.Month, result)
	}
	return result
}

func (e Event) Validate() error {
	if e.Name == "" && (e.Names == nil) {
		return ErrMissingNameOrNames
	}
	if e.Name != "" && e.Names != nil {
		return ErrNameOrNames
	}
	if e.Name == "" && (e.Names[0] == "" || e.Names[1] == "") {
		return ErrNamesArePair
	}
	return e.Type.Validate()
}
