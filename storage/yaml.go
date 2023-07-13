package storage

import (
	"os"
	"strings"
	"time"

	"git.sr.ht/~tymek/przypominajka/models"
	"gopkg.in/yaml.v3"
)

type YAML map[time.Month]map[int]models.Events

var _ Interface = YAML{}

func NewYAML(path string) (YAML, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var y YAML
	if err := yaml.Unmarshal(b, &y); err != nil {
		return nil, err
	}
	for m, month := range y {
		for d, day := range month {
			for i, e := range day {
				y[m][d][i].Month = m
				y[m][d][i].Day = d
				if err := e.Validate(); err != nil {
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
