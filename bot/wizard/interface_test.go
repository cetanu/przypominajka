package wizard

import "testing"

func TestParseCallbackData(t *testing.T) {
	tests := []struct {
		s        string
		id       string
		static   []string
		expected string
		err      error
	}{
		{
			"add:x:month:1",
			"x",
			[]string{"month"},
			"1",
			nil,
		},
		{
			"add:x:1",
			"x",
			[]string{},
			"1",
			nil,
		},
		{
			"add:x",
			"x",
			[]string{"month"},
			"",
			ErrInvalidCallbackData,
		},
		{
			"add:x:month",
			"x",
			[]string{"month"},
			"",
			ErrInvalidCallbackData,
		},
		{
			"add:x:month:",
			"x",
			[]string{"month"},
			"",
			nil,
		},
		{
			"add:asdf:month:",
			"x",
			[]string{"month"},
			"",
			ErrInvalidCallbackData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			result, err := parseCallbackData(tt.s, &Add{id: tt.id}, tt.static...)
			if result != tt.expected || err != tt.err {
				t.Errorf("got value='%s', err='%v'; expected value='%s', err='%v'", result, err, tt.expected, tt.err)
			}
		})
	}
}
