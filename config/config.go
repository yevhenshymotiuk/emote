package config

import (
	"log"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
)

// Config represents config data that is stored in the file
type Config struct {
	Dest     string            `toml:"dest"`
	Emoticon map[string]string `toml:"emoticon"`
}

type configAlreadyExistsError struct{}

func (e *configAlreadyExistsError) Error() string {
	return "config file already exists"
}

// Load unmarshals config data to Config object
func Load(v *viper.Viper) (*Config, error) {
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if _, err = Create(v, Config{Dest: "clipboard"}); err != nil {
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

// Create creates config file using defaults from Config object
func Create(v *viper.Viper, c Config) (*Config, error) {
	if err := v.ReadInConfig(); err == nil {
		return nil, &configAlreadyExistsError{}
	}

	f, err := os.Create(v.ConfigFileUsed())
	if err != nil {
		return nil, err
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

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

// File returns path to the config file
func File(v *viper.Viper) string {
	return v.ConfigFileUsed()
}
