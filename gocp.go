package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// set subcommand
	hello := flag.NewFlagSet("hello", flag.ExitOnError)

	// init
	// make directory and template files (default language is C++)
	init := flag.NewFlagSet("init", flag.ExitOnError)
	/// parse input contest No and language (if necessarily)
	/// make directory
	/// scrape contest page
	//// collect task list
	///// make directory by task name
	////// scrape task sentence and print it into file
	////// scrape test case input and output, and print them into files
	////// make template files

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
	case "hello":
		hello.Parse(os.Args[1:])
		fmt.Println("hello")
	case "init":
		init.Parse(os.Args[0:])
	case "info":
		info.Parse(os.Args[0:])
	case "test":
		test.Parse(os.Args[0:])
	case "submit":
		submit.Parse(os.Args[0:])
	default:
		flag.Usage()
	}
}
