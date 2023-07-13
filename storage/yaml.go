package storage

import (
	"os"
	"strings"
	"time"

	"git.sr.ht/~tymek/przypominajka/models"
	"gopkg.in/yaml.v3"
)

type YAML struct {
	path string
	data map[time.Month]map[int]models.Events
}

var _ Interface = (*YAML)(nil)

func NewYAML(path string) (*YAML, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	y := &YAML{path: path}
	if err := yaml.Unmarshal(b, &y.data); err != nil {
		return nil, err
	}
	for m, month := range y.data {
		for d, day := range month {
			for i, e := range day {
				y.data[m][d][i].Month = m
				y.data[m][d][i].Day = d
				if err := e.Validate(); err != nil {
					return nil, err
				}
			}
		}
	}
	return y, nil
}

func (y *YAML) String() string {
	lines := []string{}
	for m := time.January; m <= time.December; m++ {
		for d := 1; d <= 31; d++ {
			if events, ok := y.data[m][d]; ok {
				lines = append(lines, events.String())
			}
		}
	}
	return strings.Join(lines, "\n")
}

func (y *YAML) At(t time.Time) (models.Events, error) {
	return y.data[t.Month()][t.Day()], nil
}

func (y *YAML) Add(e models.Event) error {
	if _, ok := y.data[e.Month]; !ok {
		y.data[e.Month] = map[int]models.Events{}
	}
	y.data[e.Month][e.Day] = append(y.data[e.Month][e.Day], e)
	b, err := yaml.Marshal(y.data)
	if err != nil {
		return err
	}
	return os.WriteFile(y.path, b, 0o644)
}
