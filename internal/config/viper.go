package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	PostgresHost     string `mapstructure:"POSTGRES_HOST"`
	PostgresDB       string `mapstructure:"POSTGRES_DB"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
}

func LoadConfig() (config Config, err error) {
	viper.AutomaticEnv()

	err = viper.BindEnv("POSTGRES_HOST")
	if err != nil {
		return
	}

	err = viper.BindEnv("POSTGRES_PASSWORD")
	if err != nil {
		return
	}

	err = viper.BindEnv("POSTGRES_USER")
	if err != nil {
		return
	}

	err = viper.BindEnv("POSTGRES_DB")
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}

func (c Config) ConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", c.PostgresHost, c.PostgresUser, c.PostgresPassword, c.PostgresDB)
}
