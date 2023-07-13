package yaml

import (
	"errors"
	"time"

	"git.sr.ht/~tymek/przypominajka/models"
)

type YAML map[time.Month]map[int]models.Events

var _ models.Storage = YAML{}

func NewYAML() (YAML, error) {
	return nil, errors.New("not implemented")
}

func (y YAML) At(t time.Time) (models.Events, error) {
	return y[t.Month()][t.Day()], nil
}
