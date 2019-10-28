package contest

import (
	"github.com/matsuyoshi30/gocp/util"
)

func CheckTasks(contestNo string) ([]string, error) {
	alpha := "abcdefghijklmnopqrstuvwxyz"

	tasks := make([]string, 0)
	url := contestNo + "/tasks/" + contestNo
	for i := 0; i < len(alpha); i++ {
		taskURL := url + "_" + string(alpha[i])
		err := util.ValidateHeader(taskURL)
		if err != nil {
			if len(tasks) == 0 {
				return nil, err
			}
			return tasks, nil
		}
		util.LogWrite(util.SUCCESS, "Access to contest page: "+taskURL)
		tasks = append(tasks, string(alpha[i]))
	}

	return tasks, nil
}
