package main

import (
	"errors"
	"fmt"
	"net/http"
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
		return errors.New("ERROR: status code")
	}

	return nil
}

func checkTasks(contestNo string) ([]string, error) {
	alpha := "abcdefghijklmnopqrstuvwxyz"

	tasks := make([]string, 0)
	url := contestNo + "/tasks/" + contestNo
	for i := 0; i < len(alpha); i++ {
		taskURL := url + "_" + string(alpha[i])
		err := validateHeader(taskURL)
		if err != nil {
			if len(tasks) == 0 {
				return nil, err
			}
			return tasks, nil
		}
		logWrite(SUCCESS, "Access to contest page: "+taskURL)
		tasks = append(tasks, string(alpha[i]))
	}

	return tasks, nil
}

type Status int

const (
	SUCCESS Status = iota
	FAILED
)

func logWrite(st Status, str string) {
	switch st {
	case SUCCESS:
		fmt.Println("[SUCCESS]", str)
	case FAILED:
		fmt.Println("[FAILED]", str)
	default:
		fmt.Println(str)
	}
}
