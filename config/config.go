package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func ReadConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("hashcode")
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	viper.SetDefault("source-dir", cwd)
	viper.SetDefault("datasets-dir", cwd)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	return nil
}

func Token() (string, error) {
	token := strings.TrimSpace(viper.GetString("token"))
	if token == "" {
		return "", fmt.Errorf("Token is empty, please add token in config file or fill %s env variable", "HASHCODE_TOKEN")
	}
	return token, nil
}

func Datasets() (map[string]string, error) {
	datasets := viper.GetStringMapString("datasets")
	if len(datasets) == 0 {
		return nil, errors.New("No datasets defined, please add datasets in config file")
	}
	return datasets, nil
}

func SourceDir() string {
	return viper.GetString("source-dir")
}

func DatasetsDir() string {
	return viper.GetString("datasets-dir")
}
