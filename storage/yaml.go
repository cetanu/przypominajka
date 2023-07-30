package storage

import (
	"os"
	"strings"
	"time"

	"git.sr.ht/~tymek/przypominajka/models"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

type YAML struct {
	path string
	// Chat ID -> Month -> Day -> Events
	data      map[int64]map[time.Month]map[int]models.Events
	chatIDs   []int64
	chatIDMap map[int64]struct{}
}

var _ Interface = (*YAML)(nil)

type dtoYAML struct {
	Data    map[int64]map[time.Month]map[int]models.Events `yaml:"data"`
	ChatIDs []int64                                        `yaml:"chat_ids"`
}

func NewYAML(path string) (*YAML, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var dto dtoYAML
	if err := yaml.Unmarshal(b, &dto); err != nil {
		return nil, err
	}
	m := map[int64]struct{}{}
	for _, id := range dto.ChatIDs {
		m[id] = struct{}{}
	}
	y := &YAML{
		path:      path,
		data:      dto.Data,
		chatIDs:   dto.ChatIDs,
		chatIDMap: m,
	}
	for id, chat := range y.data {
		for m, month := range chat {
			for d, day := range month {
				for i, e := range day {
					y.data[id][m][d][i].Month = m
					y.data[id][m][d][i].Day = d
					if err := e.Validate(); err != nil {
						return nil, err
					}
				}
			}
		}
	}
	return y, nil
}

func (y *YAML) Format(chatID int64) string {
	lines := []string{}
	for m := time.January; m <= time.December; m++ {
		for d := 1; d <= 31; d++ {
			if events, ok := y.data[chatID][m][d]; ok {
				lines = append(lines, events.String())
			}
		}
	}
	if len(lines) == 0 {
		return "Nie ma żadnych wydarzeń"
	}
	return strings.Join(lines, "\n")
}

func (y *YAML) ChatIDs() []int64 {
	return y.chatIDs
}

func (y *YAML) IsEnabled(chatID int64) bool {
	_, ok := y.chatIDMap[chatID]
	return ok
}

func (y *YAML) At(chatID int64, t time.Time) (models.Events, error) {
	return y.data[chatID][t.Month()][t.Day()], nil
}

func (y *YAML) Add(chatID int64, e models.Event) error {
	if _, ok := y.data[chatID]; !ok {
		y.data[chatID] = map[time.Month]map[int]models.Events{}
	}
	if _, ok := y.data[chatID][e.Month]; !ok {
		y.data[chatID][e.Month] = map[int]models.Events{}
	}
	y.data[chatID][e.Month][e.Day] = append(y.data[chatID][e.Month][e.Day], e)
	return y.write()
}

func (y *YAML) Remove(chatID int64, e models.Event) error {
	if _, ok := y.data[chatID]; !ok {
		return nil
	}
	if _, ok := y.data[chatID][e.Month]; !ok {
		return nil
	}
	if _, ok := y.data[chatID][e.Month][e.Day]; !ok {
		return nil
	}
	y.data[chatID][e.Month][e.Day] = slices.DeleteFunc(y.data[chatID][e.Month][e.Day], func(other models.Event) bool {
		return cmp.Equal(e, other)
	})
	if len(y.data[chatID][e.Month][e.Day]) == 0 {
		delete(y.data[chatID][e.Month], e.Day)
	}
	if len(y.data[chatID][e.Month]) == 0 {
		delete(y.data[chatID], e.Month)
	}
	if len(y.data[chatID]) == 0 {
		delete(y.data, chatID)
	}
	return y.write()
}

func (y *YAML) write() error {
	dto := dtoYAML{
		ChatIDs: y.chatIDs,
		Data:    y.data,
	}
	b, err := yaml.Marshal(dto)
	if err != nil {
		return err
	}
	return os.WriteFile(y.path, b, 0o644)
}
