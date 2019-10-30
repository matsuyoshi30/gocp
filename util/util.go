package util

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func GetCredentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", "", err
	}
	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

func ValidateHeader(url string) error {
	resp, err := http.Head(url)
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
