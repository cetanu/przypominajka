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
		return "Nie ma żadnych wydarzeń"
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
	parts := []string{}
	if list {
		parts = append(parts, fmt.Sprintf("%02d.%02d", e.Day, e.Month), "-")
	}
	if e.Name != "" {
		parts = append(parts, e.Name)
	} else {
		parts = append(parts, e.Names[0], "i", e.Names[1])
	}
	if e.Surname != "" {
		parts = append(parts, e.Surname)
	}
	if e.Name != "" {
		parts = append(parts, "obchodzi")
	} else {
		parts = append(parts, "obchodzą")
	}
	if !list {
		parts = append(parts, "dziś")
	}
	parts = append(parts, e.Type.String())

	result := strings.Join(parts, " ")
	if !list {
		result += "!"
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
