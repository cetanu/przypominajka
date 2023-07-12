package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type year map[time.Month]map[int]events

var _ fmt.Stringer = year{}

func (y year) next() (time.Month, int, events) {
	now := time.Now()
	nextYear := now.AddDate(1, 0, 0)
	for t := now; t.Before(nextYear); t = t.AddDate(0, 0, 1) {
		if events := y.at(t); len(events) > 0 {
			return t.Month(), t.Day(), events
		}
	}
	return 0, 0, nil
}

func (y year) at(t time.Time) events {
	return y[t.Month()][t.Day()]
}

func (y year) today() events {
	return y.at(time.Now())
}

func (y year) String() string {
	lines := []string{}
	for m := time.January; m <= time.December; m++ {
		month, ok := y[m]
		if !ok {
			continue
		}
		for d := 1; d <= 31; d++ {
			day, ok := month[d]
			if !ok {
				continue
			}
			lines = append(lines, day.format(m, d))
		}
	}
	return strings.Join(lines, "\n")
}

func readYear(path string) (year, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config struct {
		January   map[int]events `yaml:"january"`
		February  map[int]events `yaml:"february"`
		March     map[int]events `yaml:"march"`
		April     map[int]events `yaml:"april"`
		May       map[int]events `yaml:"may"`
		June      map[int]events `yaml:"june"`
		July      map[int]events `yaml:"july"`
		August    map[int]events `yaml:"august"`
		September map[int]events `yaml:"september"`
		October   map[int]events `yaml:"october"`
		November  map[int]events `yaml:"november"`
		December  map[int]events `yaml:"december"`
	}
	if err := yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	return year{
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
