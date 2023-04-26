package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Jdbc struct {
		DriverClassName   string `yaml:"driverClassName"`
		Url               string `yaml:"url"`
		Username          string `yaml:"username"`
		Password          string `yaml:"password"`
		InitialSize       int    `yaml:"initialSize"`
		MaxTotal          int    `yaml:"maxTotal"`
		MaxIdle           int    `yaml:"maxIdle"`
		DefaultAutoCommit bool   `yaml:"defaultAutoCommit"`
		TestOnBorrow      bool   `yaml:"testOnBorrow"`
		ValidationQuery   string `yaml:"validationQuery"`
	} `yaml:"jdbc"`
}

func NewConfig(args []string) (*Config, error) {
	merged, err := mergeYamls(args)
	if err != nil {
		return nil, err
	}
	b, err := yaml.Marshal(merged)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(b, &config)
	return &config, err
}

func mergeYamls(args []string) (map[string]interface{}, error) {
	if len(args) == 0 {
		return nil, errors.New("=> Specify yaml path")
	}

	var err error
	var merged map[string]interface{}
	for _, arg := range args {
		data, err := readYaml(arg)
		if err != nil {
			return nil, err
		}
		if merged == nil {
			merged = data
			continue
		}
		merged, err = mergeMap(merged, data)
		if err != nil {
			return nil, err
		}
	}
	return merged, err
}

func readYaml(path string) (map[string]interface{}, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = yaml.Unmarshal(b, &data)
	return data, err
}

func mergeMap(map1, map2 map[string]interface{}) (map[string]interface{}, error) {
	var err error
	for key, val := range map2 {
		switch val.(type) {
		case map[string]interface{}:
			if _, ok := map1[key]; !ok {
				map1[key] = val
			} else {
				map1[key], err = mergeMap(map1[key].(map[string]interface{}), val.(map[string]interface{}))
				if err != nil {
					return nil, err
				}
			}
		default:
			map1[key] = val
		}
	}
	return map1, err
}

func (config *Config) ToDsn() string {
	// user:password@tcp(127.0.0.1:3306)/database_name?charset=utf8mb4&parseTime=True&loc=Local
	// user:password@tcp(127.0.0.1:3306)/database_name?characterEncoding=UTF-8&characterSetResults=UTF-8&useCursorFetch=true&defaultFetchSize=1000&autoReconnect=true&useSSL=false
	url := config.Jdbc.Url
	// url = reg(`^.*//(.*)/(.*)`).ReplaceAllString(url, "tcp($1)/$2")
	url = reg(`^.*//(.*)/(.*)`).ReplaceAllString(url, "tcp(localhost:3306)/$2")
	url = reg("\\?.*").ReplaceAllString(url, "?charset=utf8mb4&parseTime=True&loc=Local")
	return fmt.Sprintf("%s:%s@%s", config.Jdbc.Username, config.Jdbc.Password, url)
}

func reg(str string) *regexp.Regexp {
	return regexp.MustCompile(str)
}
