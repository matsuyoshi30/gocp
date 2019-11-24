package contest

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

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

	if resp.StatusCode != 200 {
		return errors.New("unexpected status code")
	}
	// if input wrong id or password, response status code is 200
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	_, err = util.Scrape(string(b), "title")
	if err != nil {
		return err
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
		err := client.ValidateHeader(baseURL + "/contests/" + taskURL)
		if err != nil {
			if len(tasks) == 0 {
				return nil, err
			}
			// ここまで取得した task を返す
			return tasks, nil
		}
		tasks = append(tasks, string(alpha[i]))
	}

	return tasks, nil
}

func GetTestCase(contestNo, taskID string) ([]string, error) {
	url := baseURL + "/contests/" + contestNo + "/tasks/" + contestNo + "_" + taskID
	util.LogWrite(util.INFO, "Task Page URL", url)
	client, err := client.NewClient()
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	// get test cases
	testcases, err := util.Scrape(string(b), "pre")
	if err != nil {
		return nil, err
	}

	return testcases, nil
}

func Submit(cookie *http.Cookie, contestNo, taskID, code string) error {
	hc, err := client.NewClient()
	if err != nil {
		return err
	}

	// make URL
	contestURL := baseURL + "/contests/" + contestNo

	submissionURL := contestURL + "/submissions/me"

	// GET cookie and csrf_token
	var token string
	if cookie.Value != "" {
		if strings.Contains(cookie.Value, "csrf_token") {
			// FIXME
			rawToken := strings.Split(strings.Split(cookie.Value, "csrf_token%3A")[1], "_TS")[0]
			rawToken = strings.ReplaceAll(rawToken, "%00", "")
			token, err = url.QueryUnescape(rawToken)
			if err != nil {
				return err
			}
		}
	}

	// http POST
	form := url.Values{}
	form.Add("data.TaskScreenName", contestNo+"_"+taskID)
	form.Add("data.LanguageId", "3003")
	form.Add("sourceCode", code)
	form.Add("csrf_token", token)

	// make request
	req, err := http.NewRequest("POST", contestURL+"/submit", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)

	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	// FIXME time sleep
	util.LogWrite(util.INFO, "Wait judging ...")
	time.Sleep(time.Second * 1) // FIXMEEEE

	for {
		req, err = http.NewRequest("GET", submissionURL, nil)
		if err != nil {
			return err
		}
		req.AddCookie(cookie)

		resp, err = hc.Do(req)
		if err != nil {
			return err
		}

		// check result
		b, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return err
		}
		res, err := util.Scrape(string(b), "tbody")
		if err != nil {
			return err
		}

		if len(res) != 0 {
			if res[0] == "AC" {
				util.LogWrite(util.SUCCESS, "PASSED!")
				break
			} else if res[0] == "Judging" || res[0] == "WJ" {
				time.Sleep(time.Second * 1) // FIXMEEEE
				continue
			} else {
				util.LogWrite(util.FAILED, "FAILED...")
				break
			}
		}
	}

	return nil
}
