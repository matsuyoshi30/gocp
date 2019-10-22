package main

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	StatusError = errors.New("error status code")
)

const (
	baseURL = "https://atcoder.jp/contests/"
)

func validateHeader(str string) error {
	resp, err := http.Head(baseURL + str)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return StatusError
	}

	return nil
}

func checkTasks(contestNo string) []string {
	alpha := "abcdefghijklmnopqrstuvwxyz"

	tasks := make([]string, 0)
	url := contestNo + "/tasks/" + contestNo
	for i := 0; i < len(alpha); i++ {
		taskURL := url + "_" + string(alpha[i])
		logWrite("[ACCESS]", taskURL)
		err := validateHeader(taskURL)
		if err != nil {
			logWrite("[INVALID] ", taskURL)
			return tasks
		}
		logWrite("[SUCCESS] Access to", taskURL)
		tasks = append(tasks, string(alpha[i]))
	}

	return tasks
}

func logWrite(str ...string) {
	fmt.Println(str)
}
