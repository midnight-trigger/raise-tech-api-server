package configs

import (
	"bytes"
	"os"

	"github.com/spf13/viper"
)

var runtimeEnv string

var localYaml = []byte(`
slack:
  channel: ""
  endpoint: ""
mysql:
  maxIdle: 10
  maxConn: 100
`)

func Init() {
	err := loadConfig()
	if err != nil {
		panic(err)
	}
}

func IsLocal() bool {
	return runtimeEnv == "local"
}

func loadConfig() (err error) {

	viper.SetConfigType("yaml")
	runtimeEnv = os.Getenv("APP_ENV")
	filename := runtimeEnv
	switch filename {
	case "local":
		err = viper.ReadConfig(bytes.NewBuffer(localYaml))
	}
	return
}
