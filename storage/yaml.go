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
	var dto yamlDTO
	if err := yaml.Unmarshal(b, &dto); err != nil {
		return nil, err
	}
	return dtoToYAML(dto)
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

type yamlDTO struct {
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

func dtoToYAML(dto yamlDTO) (YAML, error) {
	y := YAML{
		time.January:   dto.January,
		time.February:  dto.February,
		time.March:     dto.March,
		time.April:     dto.April,
		time.May:       dto.May,
		time.June:      dto.June,
		time.July:      dto.July,
		time.August:    dto.August,
		time.September: dto.September,
		time.October:   dto.October,
		time.November:  dto.November,
		time.December:  dto.December,
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
