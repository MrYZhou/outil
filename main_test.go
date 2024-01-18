package main

import (
	"fmt"
	"os"
	"testing"

	. "github.com/MrYZhou/outil/command"
	. "github.com/MrYZhou/outil/ssh"
)
func TestRun(t *testing.T) {
	Run(".","docker stats")
}

func TestConnectWithKey(t *testing.T) {
	var cli Cli
	// 使用ioutil对文件读取字符串
	filePath := "./key/larry.pem"

	// 使用 ioutil.ReadFile 读取整个文件内容
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// 将读取到的字节切片转换为字符串
	contentStr := string(contentBytes)
	cli.Host =  "47.24.11.197"
	cli.User = "root"
	cli.Password = ""

	cli.PrivateKey = contentStr
	con, err := ConnectServer(cli)
	fmt.Println(con, err)
}

func TestServer(t *testing.T){
	var cli Cli
	var host = "192.168.0.62:22"
	var user = "root"
	var password = "YH4WfLbGPasRLVhs"
	cli.Host =  host
	cli.User = user
	cli.Password = password
	con, err := ConnectServer(cli)
	fmt.Println(con, err)
}
