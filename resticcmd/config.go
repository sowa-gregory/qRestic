package resticcmd

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Name       string
	Repository string
	Password   string
	Sources    []string
	Excludes   []string
}

var configs []Config
var configLoaded bool = false
var selectedConfig int = 0

func ReadConfig(fileName string) ([]string, error) {
	dat, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("couldn't read config file:%s", fileName)
	}
	if err = json.Unmarshal(dat, &configs); err != nil {
		return nil, fmt.Errorf("couldn't unmarshall config file:%s", fileName)
	}
	configLoaded = true

	var configNames []string
	for _, conf := range configs {
		configNames = append(configNames, conf.Name)
	}
	return configNames, nil
}

func SelectConfig(index int) {
	if index < 0 || index >= len(configs) {
		panic("invalid config index")
	}
	selectedConfig = index
}
