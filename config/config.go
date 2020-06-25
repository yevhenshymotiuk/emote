package config

import (
	"os"
	"path"

	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
)

type Config struct {
	Dest     string            `toml:"dest"`
	Emoticon map[string]string `toml:"emoticon"`
}

var (
	configPath string = os.Getenv("HOME")
	configName string = "emote"
	configFile string = path.Join(configPath, configName+".toml")
)

type configAlreadyExistsError struct{}

func (e *configAlreadyExistsError) Error() string {
	return "config file already exists"
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName(configName)
	v.AddConfigPath(configPath)

	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if _, err = Create(); err != nil {
				return nil, err
			}
			if err = v.ReadInConfig(); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	c := &Config{}
	err = v.Unmarshal(c)
	return c, err
}

func Create() (*Config, error) {
	v := viper.New()

	if err := v.ReadInConfig(); err == nil {
		return nil, &configAlreadyExistsError{}
	}

	f, err := os.Create(configFile)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	c := Config{Dest: "clipboard"}

	cfgToml, err := toml.Marshal(c)
	if err != nil {
		return nil, err
	}

	_, err = f.Write(cfgToml)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func File() string {
	return configFile
}
