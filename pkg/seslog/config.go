package seslog

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"errors"
)

const DEFAULT_CONFIG_LOCATION = "/etc/seslog/seslog.toml"
const DEFAULT_CONFIG_LOCATION2 = "config/seslog.toml"

func unique(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

func LoadConfig(configFilename string) (Options, error) {
    var options Options
    configFiles := []string{configFilename, DEFAULT_CONFIG_LOCATION, DEFAULT_CONFIG_LOCATION2}
    configFiles = unique(configFiles)
    for _, configFilename = range configFiles{
        fmt.Printf("Try to load config %s\n", configFilename)
        if _, err := toml.DecodeFile(configFilename, &options); err != nil {
            continue;
        }
        fmt.Printf("Config found %s\n", configFilename)
        return options, nil
    }
    return options, errors.New("Config file not found.")
}