package qpm

import (
	"bufio"
	"net/url"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestQPM_ReadAquifer(t *testing.T) {
	c := Config{
		AquiferPath:   "./testdata/",
		AquiferRemote: &url.URL{},
		Shell:         "zsh",
	}

	a, err := ReadStratum(c, "foo")
	if err != nil {
		t.Fatalf("%+v\n", err)
	}

	want, got := stratum{}, a
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}

	t.Fail()
}

func TestQPM_Execute(t *testing.T) {
	c := Config{
		AquiferPath:   "./testdata/",
		AquiferRemote: &url.URL{},
		Shell:         "zsh",
	}

	err := Execute(c, stratum{}, "", bufio.NewWriter(os.Stdout), bufio.NewWriter(os.Stderr))
	if err != nil {
		t.Fatal(err)
	}

	t.Fail()
}
