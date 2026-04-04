package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
		Mode string `mapstructure:"mode"`
	} `mapstructure:"server"`
	DB struct {
		Driver string `mapstructure:"driver"`
		Host   string `mapstructure:"host"`
		Port   int    `mapstructure:"port"`
		User   string `mapstructure:"user"`

		Password string `mapstructure:"password"`
		DbName   string `mapstructure:"dbname"`
		Charset  string `mapstructure:"charset"`
	} `mapstructure:"db"`
	JWT struct {
		SymmetricKey string        `mapstructure:"symmetrickey"`
		Duration     time.Duration `mapstructure:"duration"`
	} `mapstructure:"jwt"`
	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
	} `mapstructure:"redis"`
	AliPay struct {
		AppID           string `mapstructure:"appid"`
		PrivateKey      string `mapstructure:"privatekey"`
		AlipayPublicKey string `mapstructure:"alipaypublickey"`
	} `mapstructure:"alipay"`
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err = viper.Unmarshal(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
