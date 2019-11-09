package subcommand

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

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

	err = contest.Login(username, password)
	if err != nil {
		util.LogWrite(util.FAILED, "Failed to login")
		return err
	}

	util.LogWrite(util.SUCCESS, "Success login")
	return nil
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
	err := client.ValidateHeader("https://atcoder.jp/contests/" + contestNo)
	if err != nil {
		return err
	}
	util.LogWrite(util.SUCCESS, "Contest Page", "https://atcoder.jp/contests/"+contestNo)

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

			testcases, err := contest.GetTestCase(contestNo, task)
			if err != nil {
				return err
			}

			for idx, testcase := range testcases {
				filename := "out"
				if idx%2 == 0 { // input
					filename = "in"
				}
				filename = filename + "_" + strconv.Itoa(idx)

				testfile := filepath.Join(p, filename)
				util.LogWrite(util.SUCCESS, "Add testcase", filepath.Join(contestNo, task, filename))
				tf, err := os.Create(testfile)
				if err != nil {
					return err
				}
				tf.Close()

				err = ioutil.WriteFile(testfile, []byte(testcase), 0644)
				if err != nil {
					return err
				}
			}
		}
		util.LogWrite(util.SUCCESS, "Make working directory", dir)
	}

	// TODO: scrape contest page
	// scrape task sentence and print it into file
	// scrape test case input and output, and print them into files

	return nil
}

func RunTest() error {
	util.LogWrite(util.SUCCESS, "Run test")

	wd, err := os.Getwd()
	if err != nil {
		util.LogWrite(util.FAILED, "Could not get working dir")
		return err
	}

	ef := "./a.out" // executable file
	if _, err := os.Stat(ef); os.IsNotExist(err) {
		util.LogWrite(util.FAILED, "Not found executable file")
		return err
	}

	i := 0
	cnt := 1
	for {
		infn := "in_" + strconv.Itoa(i)
		outfn := "out_" + strconv.Itoa(i+1)

		if _, err = os.Stat(infn); os.IsNotExist(err) { // check input file exist
			break
		}
		if _, err = os.Stat(outfn); os.IsNotExist(err) { // check output file exist
			break
		}

		// read input file
		inf, err := os.Open(filepath.Join(wd, infn))
		if err != nil {
			return err
		}
		infb, err := ioutil.ReadAll(inf)
		if err != nil {
			return err
		}
		inval := string(infb)

		// read output file
		outf, err := os.Open(filepath.Join(wd, outfn))
		if err != nil {
			return err
		}
		outfb, err := ioutil.ReadAll(outf)
		if err != nil {
			return err
		}
		outval := string(outfb)

		// execution
		cmd := exec.Command(ef)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}
		io.WriteString(stdin, inval)
		stdin.Close()
		out, err := cmd.Output()
		if err != nil {
			return err
		}
		output := strings.Trim(string(out), "\r\n")

		res := fmt.Sprintf("Case No.%d", cnt)
		comp := strings.Compare(outval, output)
		if comp != 0 {
			res = fmt.Sprintf("%s: expected %s, but got %s", res, outval, output)
			util.LogWrite(util.FAILED, res)
		} else {
			res = fmt.Sprintf("%s: PASSED! expected %s, and got %s", res, outval, output)
			util.LogWrite(util.SUCCESS, res)
		}

		i = i + 2
		cnt++
	}

	// run test

	return nil
}

func Submit() error {
	util.LogWrite(util.SUCCESS, "Submission")

	wd, err := os.Getwd()
	if err != nil {
		util.LogWrite(util.FAILED, "Could not get working dir")
		return err
	}

	sourcefile := "main.cpp" // TODO: add support optional
	if _, err := os.Stat(sourcefile); os.IsNotExist(err) {
		util.LogWrite(util.FAILED, "Not found source file")
		return err
	}

	// read sourcefile
	f, err := os.Open(filepath.Join(wd, sourcefile))
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	// check session
	ok, err := client.CheckSession(config.ConfigFile)
	if err != nil || !ok {
		util.LogWrite(util.FAILED, "Session failed")
		return err
	}

	// load cookie
	cookie, err := client.LoadCookies(config.ConfigFile)
	if err != nil {
		return err
	}

	// get taskID and contestNo
	taskDir, err := os.Open(filepath.Join(wd, ".."))
	if err != nil {
		return err
	}
	taskID := strings.Trim(strings.TrimPrefix(wd, taskDir.Name()), "/")

	contestDir, err := os.Open(filepath.Join(wd, "../.."))
	if err != nil {
		return err
	}
	contestNo := strings.TrimSuffix(strings.TrimPrefix(wd, contestDir.Name()+"/"), "/"+taskID)

	// submit
	err = contest.Submit(cookie, contestNo, taskID, string(b))
	if err != nil {
		return err
	}

	return nil
}

func Logout() error {
	// check config
	if ok := config.IsExistConfig(config.ConfigDir, config.ConfigFile); ok {
		ok, err := client.CheckSession(config.ConfigFile)
		if err != nil {
			util.LogWrite(util.FAILED, "Does not login")
			return err
		}
		if ok {
			err = os.Remove(filepath.Join(config.ConfigDir, config.ConfigFile))
			if err != nil {
				util.LogWrite(util.FAILED, "Failed to remove config file")
				return err
			}
		}
	}

	util.LogWrite(util.SUCCESS, "Success logout")
	return nil
}
