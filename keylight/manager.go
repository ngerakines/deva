package keylight

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/mdns"
)

type Info struct {
	ProductName         string   `json:"productName"`
	HardwareBoardType   int      `json:"hardwareBoardType"`
	FirmwareBuildNumber int      `json:"firmwareBuildNumber"`
	FirmwareVersion     string   `json:"firmwareVersion"`
	SerialNumber        string   `json:"serialNumber"`
	DisplayName         string   `json:"displayName"`
	Features            []string `json:"features"`
}

type Lights struct {
	NumberOfLights int     `json:"numberOfLights"`
	Lights         []State `json:"lights"`
}

type State struct {
	On          int `json:"on"`
	Brightness  int `json:"brightness"`
	Temperature int `json:"temperature"`
}

type Manager struct {
	Endpoints  []string
	State      map[string]State
	Info       map[string]Info
	Settings   map[string]Settings
	HTTPClient HTTPClient

	lock sync.Mutex
}

type Settings struct {
	PowerOnBehavior       int `json:"powerOnBehavior"`
	PowerOnBrightness     int `json:"powerOnBrightness"`
	PowerOnTemperature    int `json:"powerOnTemperature"`
	SwitchOnDurationMs    int `json:"switchOnDurationMs"`
	SwitchOffDurationMs   int `json:"switchOffDurationMs"`
	ColorChangeDurationMs int `json:"colorChangeDurationMs"`
}

func NewManager() *Manager {
	return &Manager{
		Endpoints:  make([]string, 0),
		State:      make(map[string]State),
		Info:       make(map[string]Info),
		Settings:   make(map[string]Settings),
		HTTPClient: safeHTTPClient(),
	}
}

func (m *Manager) Discover() ([]string, error) {
	entries := make(chan *mdns.ServiceEntry, 4)

	params := &mdns.QueryParam{
		Service:             "_elg._tcp",
		Domain:              "local",
		Timeout:             3 * time.Second,
		Entries:             entries,
		WantUnicastResponse: false,
	}

	go m.discoveryListener(entries)
	if err := mdns.Query(params); err != nil {
		return nil, err
	}
	time.Sleep(3 * time.Second)
	close(entries)

	return m.Endpoints, nil
}

func (m *Manager) discoveryListener(results <-chan *mdns.ServiceEntry) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for {
		select {
		case entry, ok := <-results:
			if !ok {
				return
			}
			if strings.Contains(entry.Info, "md=Elgato Key Light") {
				m.Endpoints = append(m.Endpoints, fmt.Sprintf("%s:%d", entry.AddrV4.String(), entry.Port))
			}
		}
	}
}

func (m *Manager) LoadInfo() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, e := range m.Endpoints {
		info, err := getInfo(m.HTTPClient, e)
		if err != nil {
			return err
		}
		m.Info[e] = info
	}

	return nil
}

func (m *Manager) LoadState() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, e := range m.Endpoints {
		state, err := getState(m.HTTPClient, e)
		if err != nil {
			return err
		}
		m.State[e] = state
	}

	return nil
}

func (m *Manager) LoadSettings() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, e := range m.Endpoints {
		settings, err := getSettings(m.HTTPClient, e)
		if err != nil {
			return err
		}
		m.Settings[e] = settings
	}

	return nil
}

func (m *Manager) UpdateState(endpoint string, on, brightness, temp int) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	return updateState(m.HTTPClient, endpoint, State{
		On:          on,
		Brightness:  brightness,
		Temperature: temp,
	})
}
