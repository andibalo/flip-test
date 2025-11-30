package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	AppAddress         = ":8081"
	EnvDevEnvironment  = "DEV"
	EnvProdEnvironment = "PROD"
	ServiceName        = "transaction-service"
)

type Config interface {
	AppVersion() string
	AppID() string
	AppName() string
	AppEnv() string
	AppAddress() string
}

type AppConfig struct {
	App app
}

type app struct {
	AppEnv      string
	AppVersion  string
	Name        string
	Description string
	AppUrl      string
	AppID       string
}

func InitConfig() *AppConfig {
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath("../../")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Env config file not found, using environment variables")
		} else {
			fmt.Println("Error reading config file")
			panic(err)
		}
	}

	return &AppConfig{
		App: app{
			AppEnv:      viper.GetString("APP_ENV"),
			AppVersion:  viper.GetString("APP_VERSION"),
			Name:        ServiceName,
			Description: "transaction service",
			AppUrl:      viper.GetString("APP_URL"),
			AppID:       viper.GetString("APP_ID"),
		},
	}
}

func getRequiredString(key string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	}

	panic(fmt.Errorf("KEY %s IS MISSING", key))
}

func (c *AppConfig) AppVersion() string {
	return c.App.AppVersion
}

func (c *AppConfig) AppID() string {
	return c.App.AppID
}

func (c *AppConfig) AppName() string {
	return c.App.Name
}

func (c *AppConfig) AppEnv() string {
	return c.App.AppEnv
}

func (c *AppConfig) AppAddress() string {
	return AppAddress
}
