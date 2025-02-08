package config

import "github.com/spf13/viper"

type Config struct {
	AppPort                 int    `mapstructure:"APP_PORT"`
	RedisHost               string `mapstructure:"REDIS_HOST"`
	RedisPort               int    `mapstructure:"REDIS_PORT"`
	RedisPassword           string `mapstructure:"REDIS_PASSWORD"`
	RedisDB                 int    `mapstructure:"REDIS_DB"`
	RateMaxRequestsByIP     int    `mapstructure:"RATE_MAX_REQUESTS_BY_IP"`
	RateMaxRequestsByToken  int    `mapstructure:"RATE_MAX_REQUESTS_BY_TOKEN"`
	RatePeriodWindowSeconds int    `mapstructure:"RATE_PERIOD_WINDOW_SECONDS"`
}

func Load(path string) (*Config, error) {
	var c *Config

	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	//viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		panic(err)
	}

	return c, nil
}
