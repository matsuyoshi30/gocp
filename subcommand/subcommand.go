package subcommand

import (
	"os"
	"path/filepath"

	"github.com/matsuyoshi30/gocp/client"
	"github.com/matsuyoshi30/gocp/config"
	"github.com/matsuyoshi30/gocp/contest"
	"github.com/matsuyoshi30/gocp/util"
)

func Login() error {
	// check config
	if ok := config.IsExistConfig(config.ConfigDir, config.ConfigFile); ok {
		ok, err := client.CheckSession(config.ConfigFile)
		if err != nil {
			return err
		}
		if ok { // already login
			return nil
		}
	}

	// login and create config to save cookie
	username, password, err := util.GetCredentials()
	if err != nil {
		return err
	}

	return contest.Login(username, password)
}

func Session() (string, error) {
	if ok := config.IsExistConfig(config.ConfigDir, config.ConfigFile); !ok {
		return "NG", nil
	}

	ok, err := client.CheckSession(config.ConfigFile)
	if err != nil {
		return "", err
	}
	if ok {
		return "OK", nil
	}
	return "NG", nil
}

func Prepare(contestNo string) error {
	// parse input contestNo
	err := util.ValidateHeader(contestNo)
	if err != nil {
		return err
	}
	util.LogWrite(util.SUCCESS, "Access to contest page: "+contestNo)

	// make working directory
	wd, _ := os.Getwd()
	dir := filepath.Join(wd, contestNo)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// collect task list
		tasks, err := contest.CheckTasks(contestNo)
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
			// TODO: choose template file language (default: C++)
			template := filepath.Join(p, "main.cpp")
			f, err := os.Create(template)
			if err != nil {
				return err
			}
			defer f.Close()
		}
		util.LogWrite(util.SUCCESS, "Make working directory")
	}

	// TODO: scrape contest page
	// scrape task sentence and print it into file
	// scrape test case input and output, and print them into files

	return nil
}
