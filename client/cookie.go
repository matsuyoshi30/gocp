package client

import (
	"encoding/json"
	"net/http"

	"github.com/matsuyoshi30/gocp/config"
)

func SaveCookie(filename string, cookie *http.Cookie) error {
	data, err := json.Marshal(cookie)
	if err != nil {
		return err
	}

	c := config.NewConfig()
	return c.WriteConfig(filename, data)
}

func LoadCookies(filename string) (*http.Cookie, error) {
	c := config.NewConfig()
	data, err := c.ReadConfig(filename)
	if err != nil {
		return nil, err
	}

	var cookie *http.Cookie
	err = json.Unmarshal(data, &cookie)
	if err != nil {
		return nil, err
	}

	return cookie, nil
}
