package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Config contains the full details returned from the configuration server.
type Config struct {
	configs     map[string]string
	FullConfigs *ConfigResponse
}

// Get is a convenience function returning the highest priority property
// returned from the configuration server.
func (c *Config) Get(property string) string {
	return c.configs[property]
}

// ConfigResponse is the full body response returned from the Spring Cloud config server
type ConfigResponse struct {
	Name            string            `json:"name"`
	Profiles        []string          `json:"profiles"`
	Label           string            `json:"label"`
	Version         string            `json:"version"`
	State           string            `json:"state"`
	PropertySources []*PropertySource `json:"propertySources"`
}

// PropertySource represents the details contained in a single file from the config server
type PropertySource struct {
	Name   string            `json:"name"`
	Source map[string]string `json:"source"`
}

// Load requests the details from the configuration server and places them in a Config object
func Load(url, application, branch string, profiles ...string) (*Config, error) {
	profList := getProfileList(profiles)
	resp, err := http.Get(fmt.Sprintf("%s/%s/%s/%s", url, application, profList, branch))
	if err != nil {
		return nil, err
	}
	return parseResponse(resp)
}

func getProfileList(profiles []string) string {
	profList := ""
	for i, profile := range profiles {
		if i > 0 {
			profList += ","
		}
		profList += profile
	}
	return profList
}

func parseResponse(resp *http.Response) (*Config, error) {
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	configResp := &ConfigResponse{}
	if err := json.Unmarshal(bytes, configResp); err != nil {
		return nil, err
	}

	configs := make(map[string]string)
	cnt := len(configResp.PropertySources)
	for i := cnt - 1; i >= 0; i-- {
		ps := configResp.PropertySources[i]
		for k, v := range ps.Source {
			configs[k] = v
		}
	}

	return &Config{configs: configs}, nil
}
