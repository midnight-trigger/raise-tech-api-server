package configs

import (
	"os"

	"github.com/spf13/viper"
)

var runtimeEnv string

func Init(runtime string) {
	err := loadConfig(runtime)
	if err != nil {
		panic(err)
	}
}

func IsLocal() bool {
	return runtimeEnv == "local"
}

func loadConfig(runtime string) error {
	runtimeEnv = os.Getenv("APP_ENV")

	filename := runtimeEnv
	viper.SetConfigName(filename)
	viper.SetConfigType("yaml")
	switch runtime {
	case "":
		viper.AddConfigPath("./configs/environment")
	case "test":
		viper.AddConfigPath("../../configs/environment")
	}
	err := viper.ReadInConfig()
	return err
}
