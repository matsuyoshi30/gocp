package util

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	baseURL = "https://atcoder.jp/contests/"
)

func ValidateHeader(str string) error {
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

type Status int

const (
	SUCCESS Status = iota
	FAILED
)

func LogWrite(st Status, str string) {
	switch st {
	case SUCCESS:
		fmt.Println("[SUCCESS]", str)
	case FAILED:
		fmt.Println("[FAILED]", str)
	default:
		fmt.Println(str)
	}
}
