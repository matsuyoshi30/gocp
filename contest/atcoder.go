package contest

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/matsuyoshi30/gocp/client"
	"github.com/matsuyoshi30/gocp/config"
	"github.com/matsuyoshi30/gocp/util"
)

const baseURL = "https://atcoder.jp"

func Login(username, password string) error {
	targetURL := baseURL + "/login"

	hc, err := client.NewClient()
	if err != nil {
		return err
	}

	/// check argument url
	resp, err := hc.Get(targetURL)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("invalid URL")
	}

	// GET cookie and csrf_token
	var token string
	for _, c := range resp.Cookies() {
		if c.Value != "" {
			if strings.Contains(c.Value, "csrf_token") {
				// FIXME
				rawToken := strings.Split(strings.Split(c.Value, "csrf_token%3A")[1], "_TS")[0]
				rawToken = strings.ReplaceAll(rawToken, "%00", "")
				token, err = url.QueryUnescape(rawToken)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	// http POST
	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)
	if token != "" {
		form.Add("csrf_token", token)
	}

	req, err := http.NewRequest("POST", targetURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err = hc.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("invalid login POST")
	}

	// check response and save cookie
	for _, c := range resp.Cookies() {
		if c.Value != "" {
			err = client.SaveCookie(config.ConfigFile, c)
			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func CheckTasks(contestNo string) ([]string, error) {
	alpha := "abcdefghijklmnopqrstuvwxyz"

	tasks := make([]string, 0)
	url := contestNo + "/tasks/" + contestNo
	for i := 0; i < len(alpha); i++ {
		taskURL := url + "_" + string(alpha[i])
		err := util.ValidateHeader(baseURL + "/contests/" + taskURL)
		if err != nil {
			if len(tasks) == 0 {
				return nil, err
			}
			// ここまで取得した task を返す
			return tasks, nil
		}
		util.LogWrite(util.SUCCESS, "Access to contest page: "+taskURL)
		tasks = append(tasks, string(alpha[i]))
	}

	return tasks, nil
}
