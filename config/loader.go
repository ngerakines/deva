package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"
)

type ElgatoConfig struct {
	KeyLight []string `json:"keylight,omitempty"`
}

type NanoleafConfig struct {
	Aurora []string `json:"aurora,omitempty"`
}

type Config struct {
	Elgato   ElgatoConfig   `json:"elgato,omitempty"`
	Nanoleaf NanoleafConfig `json:"nanoleaf,omitempty"`
}

func Load(loc string) (*Config, error) {
	info, err := os.Stat(loc)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return LoadFiles(loc)
	}
	return LoadFile(loc)
}

func LoadFiles(loc string) (*Config, error) {
	c := &Config{}

	walkErr := filepath.Walk(loc, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".json") {
			return nil
		}
		if err != nil {
			return err
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		localConfig := &Config{}
		if err := json.Unmarshal(data, localConfig); err != nil {
			return err
		}

		if err := mergo.Merge(c, localConfig); err != nil {
			return err
		}

		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}

	return c, nil
}

func LoadFile(loc string) (*Config, error) {
	data, err := ioutil.ReadFile(loc)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	if err := json.Unmarshal(data, c); err != nil {
		return nil, err
	}
	return c, nil
}
