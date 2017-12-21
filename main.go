package main

import (
	"fmt"

	"github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Defaults for Kingpin flags
const (
	defaultMaxResults = "100"
)

// Icons
var (
	iconDefault = &aw.Icon{Value: "icon.png"}
	iconUpdate  = &aw.Icon{Value: "icons/update-available.png"}
)

var (
	// Kingpin and script options
	app *kingpin.Application

	// Application commands
	searchCmd *kingpin.CmdClause
	updateCmd *kingpin.CmdClause
	testCmd   *kingpin.CmdClause

	// Script options (populated by Kingpin application)
	query string

	// Workflow stuff
	wf *aw.Workflow
)

// Sets up kingpin commands
func init() {
	wf = aw.New(update.GitHub("nikitavoloboev/alfred-learn-anything"), aw.HelpURL("https://github.com/nikitavoloboev/alfred-learn-anything/issues"))
	app = kingpin.New("learn anything", "Search Learn Anything maps.")

	// Update command
	updateCmd = app.Command("update", "Check for new workflow version.")
	searchCmd = app.Command("search", "Search Learn Anything maps.")

	for _, cmd := range []*kingpin.CmdClause{
		searchCmd,
	} {
		cmd.Flag("query", "Search query.").Short('q').StringVar(&query)
	}
}

func run() {
	var err error

	cmd, err := app.Parse(wf.Args())
	if err != nil {
		wf.FatalError(err)
	}

	switch cmd {
	case searchCmd.FullCommand():
		err = doSearch()
	default:
		err = fmt.Errorf("Uknown command: %s", cmd)
	}

	// TODO: Fix error with update
	// Check for update
	// if err == nil && cmd != updateCmd.FullCommand() {
	// 	err = checkForUpdate()
	// }

	if err != nil {
		wf.FatalError(err)
	}
}

func doTest() error {
	return nil
}

// main wraps run() (actual entry point) to catch errors
func main() {
	wf.Run(run)
}
