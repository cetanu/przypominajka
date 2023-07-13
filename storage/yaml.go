package storage

import (
	"errors"
	"os"
	"strings"
	"time"

	"git.sr.ht/~tymek/przypominajka/models"
	"gopkg.in/yaml.v3"
)

var (
	errMissingNameOrNames = errors.New("'name' or 'names' must be provided")
	errNameOrNames        = errors.New("'name' is mutually exclusive with 'names'")
	errNamesArePair       = errors.New("'names' must have two elements")
	errInvalidEventType   = errors.New("invalid event type")
)

type YAML map[time.Month]map[int]models.Events

var _ Interface = YAML{}

func NewYAML(path string) (YAML, error) {
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
	y := YAML{
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
			for i, e := range day {
				y[m][d][i].Month = m
				y[m][d][i].Day = d
				if err := validate(e); err != nil {
					return nil, err
				}
			}
		}
	}
	return y, nil
}

func (y YAML) String() string {
	lines := []string{}
	for m := time.January; m <= time.December; m++ {
		for d := 1; d <= 31; d++ {
			if events, ok := y[m][d]; ok {
				lines = append(lines, events.String())
			}
		}
	}
	return strings.Join(lines, "\n")
}

func (y YAML) At(t time.Time) (models.Events, error) {
	return y[t.Month()][t.Day()], nil
}

func validate(e models.Event) error {
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
	case models.Birthday, models.Nameday, models.Wedding:
	default:
		return errInvalidEventType
	}
	return nil
}
