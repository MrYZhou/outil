package command

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

/*
执行命令

第一个参数是执行的目录,第二个参数是命令
*/
func Run(content string, direct string) error {
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
		return err
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("Command failed: ", err)
		return err
	}
	return nil
}

/*
执行命令

content 	执行的目录

direct  	执行的命令

validate	执行结果要包含这个校验字符
*/
func RunWithValidated(content string, direct string, validate string) error {
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

	// 用于存储命令的输出结果
	var outputBytes []byte
	// 创建一个管道来捕获命令的标准输出
	outputPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Failed to create stdout pipe: %v", err)
	}
	defer outputPipe.Close()

	var errBytes []byte
	// 创建一个管道来捕获命令的错误输出
	errorPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("Failed to create stderr pipe: %v", err)
	}
	defer errorPipe.Close()

	// 执行命令，并等待命令完成
	if err := cmd.Start(); err != nil {
		fmt.Println("Failed to start command: ", err)
		return err
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println("Command failed: ", err)
		return err
	}

	// 读取标准输出内容
	outputBytes, _ = io.ReadAll(outputPipe)
	// 读取错误输出内容（不在此处理，但保持打开以避免阻塞）
	errBytes, _ = io.ReadAll(errorPipe)
	
	// 命令执行不正常
	if len(errBytes) > 0 {
		return fmt.Errorf("错误信息:'%s'",errBytes)
	}

	// 命令正常执行完成,在检查输出结果是否包含校验字符
	outputStr := string(outputBytes)
	if !strings.Contains(outputStr, validate) {
		return fmt.Errorf("Validation failed: The output does not contain '%s'", validate)
	}
	
	return nil
}
