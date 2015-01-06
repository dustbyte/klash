package klash_test

import (
	"testing"

	"github.com/mota/klash"
)

func TestHelpTag(t *testing.T) {
	type TestStruct struct {
		Version bool `klash-help:"Print version and exit"`
	}
	ts := TestStruct{false}
	params, err := klash.NewParams(&ts)

	if err != nil {
		t.Error(err)
		return
	}

	if params.Mapping["version"].Help != "Print version and exit" {
		t.Error("Help message is supposed to be filled")
	}
}

func TestAliasTag(t *testing.T) {
	type TestStruct struct {
		Version bool `klash-alias:"v"`
	}
	ts := TestStruct{false}
	params, err := klash.NewParams(&ts)

	if err != nil {
		t.Error(err)
		return
	}

	if params.Mapping["version"].Alias != "v" {
		t.Error("An alias is supposed to be filled")
	}
}

func TestBadAliasTag(t *testing.T) {
	type TestStruct struct {
		Version bool `klash-alias:"v"`
		Debug   bool `klash-alias:"version"`
	}
	ts := TestStruct{false, false}
	_, err := klash.NewParams(&ts)

	if err == nil {
		t.Error("Cannot have a same/alias name for two parameters")
	}
}
