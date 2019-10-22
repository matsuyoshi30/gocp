package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func prepare(contestNo string) error {
	// parse input contest No (and language (if necessarily))
	err := validateHeader(contestNo)
	if err != nil {
		return err
	}
	logWrite(SUCCESS, "Access to contest page: "+contestNo)

	// make working directory
	wd, _ := os.Getwd()
	dir := filepath.Join(wd, contestNo)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// collect task list
		tasks, err := checkTasks(contestNo)
		if err != nil {
			return err
		}
		for _, task := range tasks {
			p := filepath.Join(dir, task)
			err = os.MkdirAll(p, os.ModePerm)
			if err != nil {
				return err
			}
			// make template files
			template := filepath.Join(p, "main.cpp")
			f, err := os.Create(template)
			if err != nil {
				return err
			}
			defer f.Close()
		}
		logWrite(SUCCESS, "Make working directory")
	}

	// TODO: scrape contest page
	// scrape task sentence and print it into file
	// scrape test case input and output, and print them into files

	return nil
}

func main() {
	// set subcommand

	// init
	// make directory and template files (default language is C++)
	prepareCommand := flag.NewFlagSet("prepare", flag.ExitOnError)

	// info
	// show contest info
	info := flag.NewFlagSet("info", flag.ExitOnError)
	/// check directory and file
	/// scrape contest No page
	/// show problem statement

	// test
	// run test cases
	test := flag.NewFlagSet("test", flag.ExitOnError)
	/// check directory and file
	/// execute program compiled by user

	// submit
	// submit file
	submit := flag.NewFlagSet("submit", flag.ExitOnError)
	/// check directory and file
	/// check session (need to login)
	/// read template file
	/// access contest submit page
	/// input into submit page
	/// execute submit
	/// read result
	/// print result

	if len(os.Args) < 2 {
		flag.Usage()
		return
	}

	switch os.Args[1] {
	case "prepare":
		prepareCommand.Parse(os.Args[2:])
	case "info":
		info.Parse(os.Args[0:])
	case "test":
		test.Parse(os.Args[0:])
	case "submit":
		submit.Parse(os.Args[0:])
	default:
		flag.Usage()
		return
	}

	if prepareCommand.Parsed() {
		if len(prepareCommand.Args()) != 1 {
			flag.Usage()
			return
		}
		err := prepare(prepareCommand.Arg(0))
		if err != nil {
			fmt.Println(err)
		}
	}
}
