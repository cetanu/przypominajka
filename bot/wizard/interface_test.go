package wizard

import "testing"

func TestParseCallbackData(t *testing.T) {
	tests := []struct {
		s        string
		static   []string
		expected string
		err      error
	}{
		{
			"add:month:1",
			[]string{"month"},
			"1",
			nil,
		},
		{
			"add:1",
			[]string{},
			"1",
			nil,
		},
		{
			"add",
			[]string{"month"},
			"",
			ErrInvalidCallbackData,
		},
		{
			"add:month",
			[]string{"month"},
			"",
			ErrInvalidCallbackData,
		},
		{
			"add:month:",
			[]string{"month"},
			"",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			result, err := parseCallbackData(tt.s, &Add{}, tt.static...)
			if result != tt.expected || err != tt.err {
				t.Errorf("got value='%s', err='%v'; expected value='%s', err='%v'", result, err, tt.expected, tt.err)
			}
		})
	}
}
