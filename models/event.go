package models

import (
	"fmt"
	"strings"
	"time"

	"git.sr.ht/~tymek/przypominajka/format"
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
	Type    EventType  `yaml:"type"`
	Month   time.Month `yaml:"-"`
	Day     int        `yaml:"-"`
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
