package main

import (
	"os"
	"sort"

	"github.com/urfave/cli/v2"

	"github.com/ngerakines/deva/commands"
)

func main() {
	app := cli.NewApp()
	app.Name = "deva"
	app.Usage = "An assistant."
	app.Copyright = "(c) 2020 Nick Gerakines"
	app.Commands = []*cli.Command{
		&commands.KeylightDiscoveryCommand,
		&commands.ModeMeetingCommand,
		&commands.ModeNormalCommand,
		&commands.ConfigValidateCommand,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
