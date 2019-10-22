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
		tasks := checkTasks(contestNo)
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
	prepareC := flag.NewFlagSet("prepare", flag.ExitOnError)

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
		prepareC.Parse(os.Args[2:])
	case "info":
		info.Parse(os.Args[0:])
	case "test":
		test.Parse(os.Args[0:])
	case "submit":
		submit.Parse(os.Args[0:])
	default:
		flag.Usage()
	}

	if prepareC.Parsed() {
		if len(prepareC.Args()) != 1 {
			fmt.Println("ERROR")
			return
		}
		err := prepare(prepareC.Arg(0))
		if err != nil {
			fmt.Println(err)
		}
	}
}
