package resticcmd

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Repository string
	Password   string
	Sources    []string
	Excludes   []string
}

var config Config
var configLoaded bool

func ReadConfig(fileName string) {
	dat, err := os.ReadFile(fileName)
	if err != nil {
		panic(fmt.Sprintf("couldn't read config file:%s", fileName))
	}
	if err = json.Unmarshal(dat, &config); err != nil {
		panic(err)
	}
	configLoaded = true
}
