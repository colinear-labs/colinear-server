package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var SingleConfig = make(map[interface{}]interface{})

func LoadSingleConfig() {
	dat, err := os.ReadFile("./config_single.yaml")
	if err != nil {
		fmt.Println(err)
	}

	yaml.Unmarshal(([]byte)(dat), &SingleConfig)

}
