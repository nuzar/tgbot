package main

import (
	"encoding/json"
	"testing"
)

type S struct {
	int
	string
}

func Test_t(t *testing.T) {
	var s S
	s.int = 2
	s.string = "b"

	b, _ := json.Marshal(s)
	t.Error(string(b))
}
