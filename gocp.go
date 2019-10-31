package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/matsuyoshi30/gocp/subcommand"
)

func main() {
	// TODO: flag.Usage() 設定

	// login
	loginCommand := flag.NewFlagSet("login", flag.ExitOnError)
	// check session
	sessionCommand := flag.NewFlagSet("session", flag.ExitOnError)
	// make directory and template files (default language is C++)
	prepareCommand := flag.NewFlagSet("prepare", flag.ExitOnError)
	// show contest info
	info := flag.NewFlagSet("info", flag.ExitOnError)
	// run test cases
	test := flag.NewFlagSet("test", flag.ExitOnError)
	// submit file
	submitCommand := flag.NewFlagSet("submit", flag.ExitOnError)

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
	case "info":
		info.Parse(os.Args[0:])
	case "test":
		test.Parse(os.Args[0:])
	case "submit":
		submitCommand.Parse(os.Args[2:])
	default:
		flag.Usage()
		return
	}

	if loginCommand.Parsed() {
		err := subcommand.Login()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if sessionCommand.Parsed() {
		fmt.Println("Check session...")
		res, err := subcommand.Session()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)
		return
	}

	if prepareCommand.Parsed() {
		if len(prepareCommand.Args()) != 1 {
			flag.Usage()
			return
		}
		err := subcommand.Prepare(prepareCommand.Arg(0))
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}
