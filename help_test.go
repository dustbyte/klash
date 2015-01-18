package klash_test

import (
	"testing"

	"github.com/mota/klash"
)

func TestUsage(t *testing.T) {
	args := struct {
		Version bool
	}{}
	help, err := klash.NewHelp("test", "dummy", &args)
	if err != nil {
		t.Error(err)
		return
	}

	usage := help.Usage()
	if usage != "Usage: test [-h] [--version] ARGS..." {
		t.Errorf("Improper usage format: %s", usage)
		return
	}
}

func TestDetails(t *testing.T) {
	args := struct {
		Version bool
	}{}
	help, err := klash.NewHelp("test", "dummy", &args)
	if err != nil {
		t.Error(err)
		return
	}

	details := help.Details()
	comp := "\t-h, --help=false        Show this help\n\t--version=false"
	if details != comp {
		t.Error("Bad details output")
	}
}

func TestCommands(t *testing.T) {
	args := struct {
		Version bool
	}{}
	help, err := klash.NewHelp("test", "dummy", &args)
	if err != nil {
		t.Error(err)
		return
	}

	sub := struct {
		Debug bool
	}{}
	subhelp, err := help.AddCommand("subcommand", "really dummy", &sub)
	if err != nil {
		t.Error(err)
		return
	}

	if subhelp.Name != "test subcommand" {
		t.Errorf("Improper subcommand usage name: %s", subhelp.Name)
	}

	comp := "\tsubcommand        really dummy"
	if help.Commands() != comp {
		t.Error("Bad commands output")
	}
}
