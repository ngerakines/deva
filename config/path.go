package config

import (
	"fmt"
	"os"
	"os/user"
)

func FirstConfig(loc string) (string, error) {
	if len(loc) > 0 {
		_, err := os.Stat(loc)
		if err != nil {
			return "", err
		}
		return loc, nil
	}

	var locations []string
	locations = append(locations, "./deva.json")

	u, err := user.Current()
	if err == nil {
		locations = append(locations, fmt.Sprintf("%s/.deva.json", u.HomeDir))
		locations = append(locations, fmt.Sprintf("%s/.deva/", u.HomeDir))
		locations = append(locations, fmt.Sprintf("%s/.config/deva/", u.HomeDir))
	}
	locations = append(locations, "/etc/deva.json")
	locations = append(locations, "/etc/deva/")

	for _, path := range locations {
		_, err := os.Stat(loc)
		if err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("no configuration found")
}
