package client

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

var RedirectAttemptedError = errors.New("redirect")

func NewClient() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 3 * time.Second,
	}

	return client, nil
}

func CheckSession(filename string) (bool, error) { // true ... already login
	cookie, err := LoadCookies(filename)
	if err != nil {
		return false, err
	}

	hc, err := NewClient()
	if err != nil {
		return false, err
	}
	hc.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return RedirectAttemptedError
	}

	// make request and add cookie
	req, err := http.NewRequest("GET", "https://atcoder.jp/contests/abc001/submit", nil)
	req.AddCookie(cookie)

	// GET submit page
	// check redirect (if redirected, not login yet)
	resp, err := hc.Do(req)
	if urlError, ok := err.(*url.Error); ok && urlError.Err == RedirectAttemptedError {
		// リダイレクトされる => 未ログイン
		return false, nil
	}
	defer resp.Body.Close()

	return true, nil
}
