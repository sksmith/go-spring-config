package springconfig

import (
	"testing"
)

func TestGet(t *testing.T) {
	configs, err := LoadWithCreds("", "go-micro-template", "master", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	val := configs.Get("app.log.level")
	if val != "debug" {
		t.Fatalf("got=%v want=%v", val, "debug")
	}
}

var data = `
top: topvalue
second:
  secondfirstchild: firstval
  secondsecondhild:
    grandchild: [a, b]
  secondthirdchild: thirdval
third:
  thirdchild: [3, 4]
fourth: 7
fifth: 8
`

func TestFlatten(t *testing.T) {
	want := map[string]interface{}{
		"top":                                 "topvalue",
		"second.secondfirstchild":             "firstval",
		"second.secondsecondchild.grandchild": "a,b",
		"second.secondthirdchild":             "thirdval",
		"third.thirdchild":                    "3,4",
		"fourth":                              "7",
		"fifth":                               "8",
	}

	got, err := flatten([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	if len(want) != len(got) {
		t.Logf("incorrect length want=%d got=%d", len(want), len(got))
		t.Fail()
	}

	for k, v := range want {
		if v != got[k] {
			t.Logf("want k=%s v=%s got k=%s v=%s", k, v, k, got[k])
		}
	}
}
