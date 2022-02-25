package app_test

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"golang.org/x/net/html"
)

var update = flag.Bool("update", false, "update golden files")

func TestGolden(t *testing.T) {
	cases := []struct {
		Name string
		ana.State
	}{
		{
			Name:  "zero_state",
			State: ana.State{},
		},
	}

	names := make(map[string]bool, len(cases))
	for _, c := range cases {
		if _, ok := names[c.Name]; ok {
			t.Errorf("duplicate case name: %s", c.Name)
		}
		names[c.Name] = true

		t.Run(c.Name, func(t *testing.T) {
			golden := filepath.Join("testdata", c.Name+".golden.html")

			var b bytes.Buffer
			if err := html.Render(&b, ana.View(c.State)); err != nil {
				t.Fatalf("html.Render error: %v", err)
			}

			if *update {
				if err := ioutil.WriteFile(golden, b.Bytes(), 0644); err != nil {
					t.Fatalf("ioutil.WriteFile error: %v", err)
				}
			}

			expected, err := ioutil.ReadFile(golden)
			if err != nil {
				t.Fatalf("ioutil.ReadFile error: %v, perhaps go test -update will fix", err)
			}

			if got, want := b.Bytes(), expected; !bytes.Equal(got, want) {
				gotGolden := fmt.Sprintf("got.%s.golden.html", c.Name)
				if err := ioutil.WriteFile(gotGolden, got, 0644); err != nil {
					t.Fatalf("ioutil.WriteFile error: %v", err)
				}
				t.Fatalf("golden mismatch: see %s and %s", gotGolden, golden)
			}
		})
	}
}
