package studyrun

import (
	"os"

	"github.com/spf13/viper"
)

const _configFormat = "toml"

func ReadConfig(filename string) (map[string]any, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	viper.SetConfigType(_configFormat)
	viper.ReadConfig(file)
	return viper.AllSettings(), nil
}
