package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/ngerakines/deva/config"
	"github.com/ngerakines/deva/keylight"
)

var ModeMeetingCommand = cli.Command{
	Name:  "mode:meeting",
	Usage: "Enable meeting mode.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Usage:   "Path to the configuration file or directory to use.",
			EnvVars: []string{"DEVA_CONFIG"},
			Value:   "deva.json",
		},
		&cli.BoolFlag{
			Name:    "validate",
			Usage:   "Validate configuration and devices",
			EnvVars: []string{"DEVA_VALIDATE"},
			Value:   false,
		},
	},
	Action: modeMeetingCommandAction,
}

func modeMeetingCommandAction(cliCtx *cli.Context) error {
	validate := cliCtx.Bool("validate")

	configPath, err := config.FirstConfig(cliCtx.String("config"))
	if err != nil {
		return err
	}
	c, err := config.Load(configPath)
	if err != nil {
		return err
	}

	if validate && len(c.Elgato.KeyLight) == 0 {
		return fmt.Errorf("no elgato key light devices found in config")
	}

	m := keylight.NewManager()
	m.Endpoints = c.Elgato.KeyLight

	if err := m.LoadInfo(); err != nil {
		return err
	}
	if err := m.LoadState(); err != nil {
		return err
	}
	if err := m.LoadSettings(); err != nil {
		return err
	}

	for _, e := range m.Endpoints {
		err = m.UpdateState(e, 1, 30, m.State[e].Temperature)
		if err != nil {
			return err
		}
	}

	return nil
}
