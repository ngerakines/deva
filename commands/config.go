package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/ngerakines/deva/config"
)

var ConfigValidateCommand = cli.Command{
	Name:  "config:validate",
	Usage: "Validate configuration",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Usage:   "Path to the configuration file or directory to use.",
			EnvVars: []string{"DEVA_CONFIG"},
			Value:   "deva.json",
		},
	},
	Action: configValidateCommandAction,
}

func configValidateCommandAction(cliCtx *cli.Context) error {
	configPath, err := config.FirstConfig(cliCtx.String("config"))
	if err != nil {
		return err
	}
	fmt.Println(configPath)
	c, err := config.Load(configPath)
	if err != nil {
		return err
	}

	encoder :=  json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	return encoder.Encode(c)
}
