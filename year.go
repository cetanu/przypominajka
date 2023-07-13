package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"git.sr.ht/~tymek/przypominajka/models"
	"gopkg.in/yaml.v3"
)

type year map[time.Month]map[int]models.Events

var (
	_ fmt.Stringer   = year{}
	_ models.Storage = year{}
)

func (y year) next() (time.Month, int, models.Events) {
	events, _ := models.Next(y)
	if len(events) > 0 {
		return events[0].Month, events[0].Day, events
	}
	return 0, 0, nil
}

func (y year) At(t time.Time) (models.Events, error) {
	return y.at(t), nil
}

func (y year) at(t time.Time) models.Events {
	return y[t.Month()][t.Day()]
}

func (y year) today() models.Events {
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
			lines = append(lines, day.Format(m, d))
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
		January   map[int]models.Events `yaml:"january"`
		February  map[int]models.Events `yaml:"february"`
		March     map[int]models.Events `yaml:"march"`
		April     map[int]models.Events `yaml:"april"`
		May       map[int]models.Events `yaml:"may"`
		June      map[int]models.Events `yaml:"june"`
		July      map[int]models.Events `yaml:"july"`
		August    map[int]models.Events `yaml:"august"`
		September map[int]models.Events `yaml:"september"`
		October   map[int]models.Events `yaml:"october"`
		November  map[int]models.Events `yaml:"november"`
		December  map[int]models.Events `yaml:"december"`
	}
	if err := yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	y := year{
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
	}
	for m, month := range y {
		for d, day := range month {
			for i := range day {
				y[m][d][i].Month = m
				y[m][d][i].Day = d
			}
		}
	}
	return y, nil
}
