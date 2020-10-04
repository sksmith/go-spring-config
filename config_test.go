package main

import "testing"

func TestGet(t *testing.T) {
	configs, err := Load("http://localhost:8888", "smfg-inventory", "master", "dev")
	if err != nil {
		t.Fatal(err)
	}
	val := configs.Get("test.property")
	if val != "dev" {
		t.Fatalf("got=%s want=%s", val, "dev")
	}
}
