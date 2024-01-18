package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	. "github.com/MrYZhou/outil/command"
	. "github.com/MrYZhou/outil/ssh"
	"golang.org/x/crypto/ssh"
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

	// 解析私钥
	signer, err := ssh.ParsePrivateKey(contentBytes)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}

	// 设置客户端请求参数
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 在生产环境中请替换为安全的主机密钥验证方法
	}

	// 作为客户端连接SSH服务器
	client, err := ssh.Dial("tcp", "47.120.11.197:22", config)
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
	var host = "192.168.0.62:22"
	var user = "root"
	var password = "YH4WfLbGPasRLVhs"
	cli.Host =  host
	cli.User = user
	cli.Password = password
	con, err := ConnectServer(cli)
	fmt.Println(con, err)
}

