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
		Shell:         []string{"zsh"},
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

func TestQPM_readStratumFile(t *testing.T) {
	t.Parallel()

	// Act
	got, err := readStratumFile("./testdata/", "foo")

	// Assert
	if cmp.Diff(nil, err) != "" {
		t.Fatalf("%+v\n", err)
	}
	want := stratumFile{
		"install": {{
			OS:         []string{"linux", "darwin"},
			Shell:      []string{"zsh", "bash", "sh"},
			Dependency: []string{"echo"},
			Step: []any{
				map[string]any{"name": "Step 1", "run": "echo install for unix"},
				"sleep 2",
				"echo install for unix",
			},
		}},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestQPM_ReadStratum(t *testing.T) {
	c := Config{
		AquiferPath: "./testdata/",
	}

	// Act
	got, err := ReadStratum(c, "foo")

	// Assert
	if cmp.Diff(nil, err) != "" {
		t.Fatalf("%+v\n", err)
	}
	want := stratum{
		Name: "foo",
		Plan: plan{},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestQPM_Execute(t *testing.T) {
	c := Config{
		AquiferPath:   "./testdata/",
		AquiferRemote: &url.URL{},
		Shell:         []string{"zsh"},
	}

	err := Execute(c, stratum{}, "", bufio.NewWriter(os.Stdout), bufio.NewWriter(os.Stderr))
	if err != nil {
		t.Fatal(err)
	}

	t.Fail()
}
