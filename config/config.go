package config

import (
	"TestTask-events/pkg/db/clickhouse"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"time"
)

const Env = "APP_CONFIG"

type Data struct {
	App      string            `yaml:"app"`
	Server   Server            `yaml:"server"`
	Database clickhouse.Config `yaml:"database"`
	Workers  Workers           `yaml:"workers"`
}

type Workers struct {
	Buffer    int `yaml:"buffer"`
	Default   int `yaml:"default"`
	Analytics int `yaml:"analytics"`
	DB        int `yaml:"db"`
}

type Server struct {
	Port        int64         `yaml:"port"`
	ReadTimeOut time.Duration `yaml:"readTimeout"`
}

func Parse(pathToConfig string) (*Data, error) {
	viper.SetConfigFile(pathToConfig)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Data{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed unmarshall yml")
	}

	return config, nil
}
