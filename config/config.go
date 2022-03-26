package config

import (
	"github.com/pkg/errors"
	"github.com/sparrowganz/TestTask-events/pkg/database"
	"github.com/spf13/viper"
	"time"
)

const Env = "APP_CONFIG"

type Data struct {
	App      string          `yml:"app"`
	Server   Server          `yml:"server"`
	Database database.Config `yml:"database"`
}

type Server struct {
	Port        int64         `yml:"port"`
	ReadTimeOut time.Duration `yml:"readTimeout"`
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
