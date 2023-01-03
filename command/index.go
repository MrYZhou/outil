package command

import (
	"os/exec"
	"strings"
)

/*
执行命令

第一个参数是执行的目录,第二个参数是命令
*/
func Run(content string,direct string) {
	param := strings.Split(direct, " ")
	cmd := exec.Command(param[0], param[1:]...)
	cmd.Dir = content
	cmd.Stdout = os.Stdout
	cmd.Run()
}
