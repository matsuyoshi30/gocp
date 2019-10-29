package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	ConfigDir  = os.Getenv("HOME") + "/Library/Preferences"
	ConfigFile = "gocp.atcoder.config.json"
)

// save cookie info
type Config struct {
	Path string
}

func NewConfig() *Config {
	return &Config{Path: ConfigDir}
}

func (c *Config) CreateConfig(filename string) error {
	// check config dir
	if err := IsExistConfig(c.Path, filename); err != nil {
		_, err = os.Create(filepath.Join(c.Path, filename))
		if err != nil {
			return err
		}
		return nil
	}

	return errors.New("config file exists.")
}

func (c *Config) ReadConfig(filename string) ([]byte, error) {
	f, err := c.OpenConfig(filename)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

func (c *Config) OpenConfig(filename string) (*os.File, error) {
	return os.Open(filepath.Join(c.Path, filename))
}

func (c *Config) WriteConfig(filename string, data []byte) error {
	return ioutil.WriteFile(filepath.Join(c.Path, filename), data, 0644)
}

func IsExistConfig(dir, filename string) error {
	if _, err := os.Stat(filepath.Join(dir, filename)); err != nil {
		return err
	}

	return nil
}
