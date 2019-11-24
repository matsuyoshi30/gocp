package util

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/html"
)

func GetCredentials(input io.Reader) (string, string, error) {
	reader := bufio.NewReader(input)

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

var NotLoginError = errors.New("Does not login")

func Scrape(source, tagtype string) ([]string, error) {
	tokens := html.NewTokenizer(strings.NewReader(source))

	ret := make([]string, 0)

	for {
		tt := tokens.Next()

		switch {
		case tt == html.ErrorToken:
			err := tokens.Err()
			if err == io.EOF {
				return ret, nil
			}
			return nil, err
		case tt == html.StartTagToken:
			tagname, _ := tokens.TagName()
			if string(tagname) == tagtype {
				if tagtype == "pre" { // read contest task page
					tokentype := tokens.Next()
					if tokentype == html.TextToken {
						text := strings.Trim(string(tokens.Text()), "\r\n")
						if len(text) != 0 { // pre element is not a testcase but input format
							ret = append(ret, text)
						}
					}
				} else if tagtype == "tbody" { // read submission/me page
					var tokentype html.TokenType
					for {
						tokentype = tokens.Next()
						if tokentype == html.TextToken {
							text := strings.Trim(string(tokens.Text()), "\r\n")
							if text == "Judging" || text == "WJ" ||
								text == "AC" || text == "WA" || text == "CE" || text == "TLE" || text == "RE" {
								ret = append(ret, text)
								return ret, nil
							}
						} else if tokentype == html.StartTagToken {
							tn, hasAttr := tokens.TagName()
							if string(tn) == "td" && hasAttr {
								key, val, moreAttr := tokens.TagAttr()
								if string(key) == "colspan" && string(val) == "3" {
									ret = append(ret, "WJ")
									return ret, nil
								}
								hasAttr = moreAttr
							}
						} else if tokentype == html.ErrorToken {
							err := tokens.Err()
							if err == io.EOF {
								return ret, nil
							}
							return nil, err
						}
					}
				} else if tagtype == "title" {
					tokens.Next()
					if strings.Trim(string(tokens.Text()), "\r\n") == "ログイン - AtCoder" {
						return nil, NotLoginError
					} else {
						return nil, nil
					}
				}
			}
		default:
			continue
		}
	}

	return ret, nil
}

type Status int

const (
	SUCCESS Status = iota
	FAILED
	INFO
)

func LogWrite(st Status, str ...string) {
	out := ""
	for _, s := range str {
		out = strings.Join([]string{out, " "}, s)
	}

	switch st {
	case SUCCESS:
		fmt.Printf("[\x1b[34m%s\x1b[0m] %s\n", "SUCCESS", out)
	case FAILED:
		fmt.Printf("[\x1b[31m%s\x1b[0m] %s\n", "FAILED", out)
	case INFO:
		fmt.Printf("[\x1b[32m%s\x1b[0m] %s\n", "INFO", out)
	default:
		fmt.Printf("%s\n", out)
	}
}
