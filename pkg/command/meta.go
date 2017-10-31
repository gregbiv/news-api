package command

import (
	"bufio"
	"flag"
	"io"
	"strings"

	"github.com/gregbiv/news-api/pkg/config"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/colorstring"
)

// Meta contains the meta-options and functionality that nearly every command inherits.
type Meta struct {
	UI cli.Ui

	// Application configuration
	Config *config.Specification

	// Whether to not-colorize output
	noColor bool
}

// FlagSet returns a FlagSet with the common flags that every command implements.
func (m *Meta) FlagSet(n string) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)

	f.BoolVar(&m.noColor, "no-color", false, "")

	// Create an io.Writer that writes to our UI properly for errors.
	// This is kind of a hack, but it does the job. Basically: create
	// a pipe, use a scanner to break it into lines, and output each line
	// to the UI. Do this forever.
	errR, errW := io.Pipe()
	errScanner := bufio.NewScanner(errR)
	go func() {
		for errScanner.Scan() {
			m.UI.Error(errScanner.Text())
		}
	}()
	f.SetOutput(errW)

	return f
}

// Colorize colorizes your strings, giving you the ability to customize
func (m *Meta) Colorize() *colorstring.Colorize {
	return &colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Disable: m.noColor,
		Reset:   true,
	}
}

// generalOptionsUsage returns the help string for the global options.
func generalOptionsUsage() string {
	helpText := `
  -no-color
    Disables colored command output.
`
	return strings.TrimSpace(helpText)
}
