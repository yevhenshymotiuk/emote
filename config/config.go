package config

import (
	"os"

	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
)

type Config struct {
	Dest     string            `toml:"dest"`
	Emoticon map[string]string `toml:"emoticon"`
}

type configAlreadyExistsError struct{}

func (e *configAlreadyExistsError) Error() string {
	return "config file already exists"
}

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

func Create(v *viper.Viper, c Config) (*Config, error) {
	if err := v.ReadInConfig(); err == nil {
		return nil, &configAlreadyExistsError{}
	}

	f, err := os.Create(v.ConfigFileUsed())
	defer f.Close()
	if err != nil {
		return nil, err
	}

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

func File(v *viper.Viper) string {
	return v.ConfigFileUsed()
}
