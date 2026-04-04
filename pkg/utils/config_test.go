package utils

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("../..")
	if err != nil {
		return
	}
	fmt.Printf("%s %s\n", config.Server.Mode, config.Server.Port)
}
