package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	. "github.com/MrYZhou/outil/command"
	. "github.com/MrYZhou/outil/ssh"
)
func TestRun(t *testing.T) {
	Run(".","docker stats")
}

func TestConnectWithKey(t *testing.T) {
	
	// 读取私钥文件内容
	contentBytes, err := os.ReadFile("C:/Users/lg/Desktop/project/go/outil/key/larry.pem")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var cli Cli
	cli.Host =  "47.120.11.197:22"
	cli.User = "root"
	cli.PrivateKey = contentBytes
	con, err := ConnectServer(cli)
	client := con.Client

	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer client.Close()

	// 创建会话
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	// 设置会话标准输出，并运行命令
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("cat /proc/cpuinfo"); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
}

func TestServer(t *testing.T){
	var cli Cli
	cli.Host =  "192.168.0.62:22"
	cli.User = "root"
	cli.Password = "YH4WfLbGPasRLVhs"
	con, err := ConnectServer(cli)
	fmt.Println(con, err)
}

