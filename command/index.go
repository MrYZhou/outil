package command

import (
	"os"
	"os/exec"
	"strings"
	"fmt"
)

/*
执行命令

第一个参数是执行的目录,第二个参数是命令
*/
func Run(content string, direct string) {
	// 将direct按空格分割为参数列表
	param := strings.Split(direct, " ")
	// 创建一个执行命令的命令对象cmd，并设置要执行的命令和参数
	cmd := exec.Command(param[0], param[1:]...)
	// 设置要执行命令的当前工作目录为content
	cmd.Dir = content
	// 将命令的标准输出重定向到os.Stdout
	cmd.Stdout = os.Stdout
	// 将命令的错误输出重定向到os.Stderr
	cmd.Stderr = os.Stderr

	// 执行命令，并等待命令完成
	if err := cmd.Start(); err != nil {
		fmt.Println("Failed to start command: ", err)
		return
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("Command failed: ", err)
		return
	}
}
