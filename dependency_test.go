package qpm

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestQPM_dependencies(t *testing.T) {
	t.Parallel()

	c := Config{
		AquiferPath: "./testdata/",
	}

	got := make(map[string][]string)

	// Act
	err := dependencies(c, Install, OS("linux"), "a", got)

	// Assert
	if cmp.Diff(nil, err) != "" {
		t.Errorf("%+v\n", err)
	}
	want := map[string][]string{
		"a": {"b", "c"},
		"b": {"c", "d", "e", "f"},
		"c": {"f"},
		"d": {"g"},
		"e": {},
		"f": {},
		"g": {},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
