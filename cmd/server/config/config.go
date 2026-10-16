package config

import (
	wcmd_utils "github.com/DVV-15324/witches/cmd/utils"
)

func Load() *wcmd_utils.Config {
	cfg := wcmd_utils.Load()
	return cfg
}
