package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/matsuyoshi30/gocp/util"
)

const USAGE = `NAME:
   gocp - a cli tool for competitive programming
USAGE:
   gocp
VERSION:
   0.1.0
COMMAND:
   login     Login competitive programming page
   session   Check session status (login or not login)
   prepare   Make directory, template file and get test cases for specified task
   test      Run test and compare output and expected value
   submit    Submit source code
   logout    Logout competitive programming page`

func main() {
	flag.Usage = func() {
		fmt.Println(USAGE)
	}

	loginCommand := flag.NewFlagSet("login", flag.ExitOnError)
	sessionCommand := flag.NewFlagSet("session", flag.ExitOnError)
	prepareCommand := flag.NewFlagSet("prepare", flag.ExitOnError)
	testCommand := flag.NewFlagSet("test", flag.ExitOnError)
	submitCommand := flag.NewFlagSet("submit", flag.ExitOnError)
	logoutCommand := flag.NewFlagSet("logout", flag.ExitOnError)

	if len(os.Args) < 2 {
		flag.Usage()
		return
	}

	switch os.Args[1] {
	case "login":
		loginCommand.Parse(os.Args[0:])
	case "session":
		sessionCommand.Parse(os.Args[0:])
	case "prepare":
		prepareCommand.Parse(os.Args[2:])
	case "test":
		testCommand.Parse(os.Args[2:])
	case "submit":
		submitCommand.Parse(os.Args[2:])
	case "logout":
		logoutCommand.Parse(os.Args[0:])
	default:
		flag.Usage()
		return
	}

	if loginCommand.Parsed() {
		err := Login()
		if err != nil {
			util.LogWrite(util.FAILED, err.Error())
			return
		}
	}

	if sessionCommand.Parsed() {
		util.LogWrite(util.INFO, "Checking session ...")
		res, err := Session()
		if err != nil {
			util.LogWrite(util.FAILED, err.Error())
			return
		}
		util.LogWrite(util.SUCCESS, res)
		return
	}

	if prepareCommand.Parsed() {
		if len(prepareCommand.Args()) != 1 {
			flag.Usage()
			return
		}
		err := Prepare(prepareCommand.Arg(0))
		if err != nil {
			util.LogWrite(util.FAILED, err.Error())
			return
		}
	}

	if testCommand.Parsed() {
		if len(testCommand.Args()) != 1 {
			flag.Usage()
			return
		}
		err := RunTest(testCommand.Arg(0))
		if err != nil {
			util.LogWrite(util.FAILED, err.Error())
			return
		}
	}

	if submitCommand.Parsed() {
		if len(submitCommand.Args()) != 1 {
			flag.Usage()
			return
		}
		err := Submit(submitCommand.Arg(0))
		if err != nil {
			util.LogWrite(util.FAILED, err.Error())
			return
		}
	}

	if logoutCommand.Parsed() {
		err := Logout()
		if err != nil {
			util.LogWrite(util.FAILED, err.Error())
			return
		}
	}
}
