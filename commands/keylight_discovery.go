package commands

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/ngerakines/deva/config"
	"github.com/ngerakines/deva/keylight"
)

var KeylightDiscoveryCommand = cli.Command{
	Name:  "keylight:discovery",
	Usage: "Discover elgato key light devices.",
	Flags: []cli.Flag{
	},
	Action: keylightDiscoveryCommandAction,
}

func keylightDiscoveryCommandAction(*cli.Context) error {
	// Hashicorp's mdns library calls `log.Printf`.
	log.SetOutput(ioutil.Discard)

	m := keylight.NewManager()
	endpoints, err := m.Discover()
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(&config.Config{
		Elgato: struct {
			KeyLight []string `json:"keylight,omitempty"`
		}{
			KeyLight: endpoints,
		},
	})
}
