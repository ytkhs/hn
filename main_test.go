package main

import (
	"context"
	"io"
	"testing"
)

type topCall struct {
	num    int
	asJSON bool
}

func runWithStub(t *testing.T, args ...string) []topCall {
	t.Helper()
	var calls []topCall
	orig := runTop
	runTop = func(num int, asJSON bool) error {
		calls = append(calls, topCall{num, asJSON})
		return nil
	}
	t.Cleanup(func() { runTop = orig })

	app := newApp()
	app.Writer = io.Discard
	if err := app.Run(context.Background(), append([]string{"hn"}, args...)); err != nil {
		t.Fatalf("Run(%v): %v", args, err)
	}
	return calls
}

func TestRouting(t *testing.T) {
	cases := []struct {
		args []string
		want []topCall
	}{
		{nil, []topCall{{10, false}}},
		{[]string{"30"}, []topCall{{30, false}}},
		{[]string{"--json"}, []topCall{{10, true}}},
		{[]string{"--json", "5"}, []topCall{{5, true}}},
		{[]string{"5", "--json"}, []topCall{{5, true}}},
		{[]string{"top"}, []topCall{{10, false}}},
		{[]string{"top", "--json", "3"}, []topCall{{3, true}}},
		{[]string{"--help"}, nil},
		{[]string{"help"}, nil},
	}
	for _, c := range cases {
		got := runWithStub(t, c.args...)
		if len(got) != len(c.want) {
			t.Errorf("%v: calls = %v, want %v", c.args, got, c.want)
			continue
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("%v: calls = %v, want %v", c.args, got, c.want)
			}
		}
	}
}
