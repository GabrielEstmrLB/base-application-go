package main_configs_env

import (
	main_error "baseapplicationgo/main/configs/error"
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
)

const _MSG_ENV_BEAN = "Invalid env property value."
const _MSG_ERROR_READ_ENV_FILE = "Error to read .env file."

const YML_BASE_DIRECTORY_MAIN_REFERENCE = "./zresources"

var once sync.Once
var envValues *map[string]string

func GetEnvConfigBean() *map[string]string {
	once.Do(func() {
		if envValues == nil {
			envValues = getEnvConfig()
		}
	})
	return envValues
}

func getEnvConfig() *map[string]string {

	envNames := []string{
		MP_INDICATOR_APPLICATION_PROFILE.GetDescription(),
	}

	err := godotenv.Load(YML_BASE_DIRECTORY_MAIN_REFERENCE + "/.env")
	main_error.FailOnError(err, _MSG_ERROR_READ_ENV_FILE)

	data := make(map[string]string)
	for _, value := range envNames {
		envValue := os.Getenv(value)
		if envValue == "" {
			log.Panicf("%s: %s", _MSG_ENV_BEAN, errors.New(_MSG_ENV_BEAN))
		}
		data[value] = envValue
	}
	return &data
}
