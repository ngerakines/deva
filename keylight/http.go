package keylight

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func safeHTTPClient() HTTPClient {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

func getInfo(client HTTPClient, endpoint string) (Info, error) {
	var info Info
	if err := request(client, "GET", fmt.Sprintf("http://%s/elgato/accessory-info", endpoint), nil, &info); err != nil {
		return Info{}, err
	}
	return info, nil
}

func getSettings(client HTTPClient, endpoint string) (Settings, error) {
	var settings Settings
	if err := request(client, "GET", fmt.Sprintf("http://%s/elgato/lights/settings", endpoint), nil, &settings); err != nil {
		return Settings{}, err
	}
	return settings, nil
}

func getState(client HTTPClient, endpoint string) (State, error) {
	var l Lights
	if err := request(client, "GET", fmt.Sprintf("http://%s/elgato/lights", endpoint), nil, &l); err != nil {
		return State{}, err
	}
	if l.NumberOfLights != 1 {
		return State{}, fmt.Errorf("unexpected response from device")
	}
	return l.Lights[0], nil
}

func updateState(client HTTPClient, endpoint string, state State) error {
	l := Lights{NumberOfLights: 1, Lights: []State{state}}
	data, err := json.Marshal(l)
	if err != nil {
		return err
	}

	if err := request(client, "PUT", fmt.Sprintf("http://%s/elgato/lights", endpoint), bytes.NewReader(data), nil); err != nil {
		return err
	}
	if l.NumberOfLights != 1 {
		return fmt.Errorf("unexpected response from device")
	}
	return nil
}

func request(client HTTPClient, method, url string, reader io.Reader, target interface{}) error {
	request, err := http.NewRequest(method, url, reader)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error: unexpected status code from device: %d", response.StatusCode)
	}

	defer response.Body.Close()
	if target != nil {
		return json.NewDecoder(response.Body).Decode(target)
	}
	return nil
}
