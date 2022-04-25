package configprovider

import (
	"Assignment/providers"
	"os"
	"strconv"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func NewConfigProvider() providers.ConfigProvider {
	return &Config{}
}

func (c *Config) Read() {
	err := envconfig.Process("", c)
	if err != nil {
		logrus.Fatal(err.Error())
	}
}

func (c *Config) GetServerPort() string {
	if c == nil {
		return ""
	}
	return c.Port
}

func (c *Config) GetString(key string) string {
	return os.Getenv(key)
}

func (c *Config) GetInt(key string) int {
	intVal, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		logrus.Errorf("Error getting int config %v", err)
	}
	return intVal
}

func (c *Config) GetAny(key string) interface{} {
	return os.Getenv(key)
}
