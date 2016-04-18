package config

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	conf, loadErr := LoadConfig("./test_conf.toml")
	if loadErr != nil {
		t.Error(loadErr)
	}
	t.Log(conf)
}
