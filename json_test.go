package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

type A struct {
	Name string `json:name`
}

func TestA(t *testing.T) {
	a := A{"hello"}
	s, _ := json.Marshal(a)
	fmt.Println(string(s))
}
