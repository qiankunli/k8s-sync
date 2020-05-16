package config

import (
	"fmt"
	"testing"
)

func TestReadYaml(t *testing.T) {
	config,err := ReadYaml("config.yaml")
	fmt.Println(err)
	fmt.Println(config.Debug)
	fmt.Println(config.BearerToken)
	fmt.Println(config.IntervalSeconds)
}
