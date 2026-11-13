package main

import "testing"

func TestParseNum(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"", 10},
		{"30", 30},
		{"500", 500},
		{"501", 10},
		{"0", 10},
		{"-5", 10},
		{"abc", 10},
	}
	for _, c := range cases {
		if got := parseNum(c.in); got != c.want {
			t.Errorf("parseNum(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}
