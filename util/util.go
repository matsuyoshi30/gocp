package util

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/html"
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

// scrape web page

func Scrape(source, tagtype string) ([]string, error) {
	LogWrite(SUCCESS, "Start Scraping")
	tokens := html.NewTokenizer(strings.NewReader(source))

	testcases := make([]string, 0)

	for {
		tt := tokens.Next()

		switch {
		case tt == html.ErrorToken:
			err := tokens.Err()
			if err == io.EOF {
				return testcases, nil
			}
			return nil, err
		case tt == html.StartTagToken:
			tagname, _ := tokens.TagName()
			if string(tagname) == tagtype {
				if tagtype == "pre" {
					tokentype := tokens.Next()
					if tokentype == html.TextToken {
						// TODO: trim whitespaces
						text := strings.Trim(string(tokens.Text()), "\r\n")
						if len(text) != 0 { // pre element is not a testcase but input format
							LogWrite(SUCCESS, "add testcase")
							testcases = append(testcases, text)
						}
					}
				} else if tagtype == "tbody" {
					// read submission/me page
					var tokentype html.TokenType
					for {
						tokentype = tokens.Next()
						if tokentype == html.TextToken {
							text := strings.Trim(string(tokens.Text()), "\r\n")
							if text == "AC" || text == "WA" || text == "CE" || text == "TLE" || text == "RE" {
								testcases = append(testcases, text)
								return testcases, nil
							}
						} else if tokentype == html.ErrorToken {
							err := tokens.Err()
							if err == io.EOF {
								return testcases, nil
							}
							return nil, err
						}
					}
				}
			}
		default:
			continue
		}
	}

	return testcases, nil
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
