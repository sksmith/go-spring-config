package springconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"gopkg.in/yaml.v2"
)

// Local Configs < Remote Configs < Environment Variables < Command Line Args

// Config contains the full details returned from the configuration server.
type Config struct {
	Values  map[string]string
	Details *ConfigDetails
}

// Get is a convenience function returning the highest priority property
// returned from the configuration server.
func (c *Config) Get(property string) string {
	return c.Values[property]
}

// ConfigResponse is the full body response returned from the Spring Cloud config server
type ConfigDetails struct {
	Name     string    `json:"name"`
	Profiles []string  `json:"profiles"`
	Label    string    `json:"label"`
	Version  string    `json:"version"`
	State    string    `json:"state"`
	Sources  []*Source `json:"propertySources"`
}

// PropertySource represents the details contained in a single file from the config server
type Source struct {
	Name    string            `json:"name"`
	Configs map[string]string `json:"source"`
}

// Load requests the details from the configuration server and places them in a Config object
func Load(url, application, branch string, profiles ...string) (*Config, error) {
	return loadFromConfigServer(url, application, branch, profiles...)
}

// Load requests the details from the configuration server and places them in a Config object
func loadFromConfigServer(url, application, branch string, profiles ...string) (*Config, error) {
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
	configResp := &ConfigDetails{}
	if err := json.Unmarshal(bytes, configResp); err != nil {
		return nil, err
	}

	configs := make(map[string]string)
	cnt := len(configResp.Sources)
	for i := cnt - 1; i >= 0; i-- {
		ps := configResp.Sources[i]
		for k, v := range ps.Configs {
			configs[k] = v
		}
	}

	return &Config{Values: configs, Details: configResp}, nil
}

func flatten(data []byte) (map[string]string, error) {
	yml := make(map[interface{}]interface{})

	err := yaml.Unmarshal([]byte(data), &yml)
	if err != nil {
		return nil, err
	}

	flatmap := make(map[string]string)
	fillflatmap("", yml, flatmap)

	return flatmap, nil
}

func fillflatmap(prefix string, yaml map[interface{}]interface{}, flatmap map[string]string) {
	for k, v := range yaml {
		if key, ok := k.(string); ok {
			if prefix != "" {
				key = prefix + "." + key
			}
			switch v.(type) {
			case string:
				flatmap[key] = v.(string)
			case int:
				flatmap[key] = strconv.Itoa(v.(int))
			case map[interface{}]interface{}:
				fillflatmap(key, v.(map[interface{}]interface{}), flatmap)
			case []interface{}:
				csv := ""
				for i, e := range v.([]interface{}) {
					if i != 0 {
						csv += ","
					}
					if s, ok := e.(string); ok {
						csv += s
					}
					if s, ok := e.(int); ok {
						csv += strconv.Itoa(s)
					}
				}
				flatmap[key] = csv
			}
		}
		if key, ok := k.(map[interface{}]interface{}); ok {
			fillflatmap(prefix, key, flatmap)
		}
	}
}
