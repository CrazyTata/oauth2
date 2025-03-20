package testutils

import (
	"oauth2/infrastructure/config"

	"github.com/zeromicro/go-zero/core/conf"
)

func LoadConfig(file string) (*config.Config, error) {
	var config config.Config
	if err := conf.Load(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
