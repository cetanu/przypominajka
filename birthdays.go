package main

import (
	_ "embed"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

//go:embed config.yaml
var bdayBytes []byte

type (
	birthdays map[time.Month]month
	month     map[int][]info
	info      struct {
		Name string    `yaml:"name"`
		Type eventType `yaml:"type"`
	}
)

const (
	birthday eventType = "birthday"
	nameday  eventType = "nameday"
)

type eventType string

func (e *eventType) UnmarshalYAML(value *yaml.Node) error {
	switch value.Value {
	case "birthday", "urodziny":
		*e = birthday
	case "nameday", "imieniny":
		*e = nameday
	default:
		return fmt.Errorf("invalid type: %s", value.Value)
	}
	return nil
}

func (b birthdays) at(t time.Time) []info {
	return b[t.Month()][t.Day()]
}

func (b birthdays) today() []info {
	return b.at(time.Now())
}

func readBirthdays() (birthdays, error) {
	var config struct {
		January   month `yaml:"january"`
		February  month `yaml:"february"`
		March     month `yaml:"march"`
		April     month `yaml:"april"`
		May       month `yaml:"may"`
		June      month `yaml:"june"`
		July      month `yaml:"july"`
		August    month `yaml:"august"`
		September month `yaml:"september"`
		October   month `yaml:"october"`
		November  month `yaml:"november"`
		December  month `yaml:"december"`
	}
	if err := yaml.Unmarshal(bdayBytes, &config); err != nil {
		return nil, err
	}
	return birthdays{
		time.January:   config.January,
		time.February:  config.February,
		time.March:     config.March,
		time.April:     config.April,
		time.May:       config.May,
		time.June:      config.June,
		time.July:      config.July,
		time.August:    config.August,
		time.September: config.September,
		time.October:   config.October,
		time.November:  config.November,
		time.December:  config.December,
	}, nil
}
