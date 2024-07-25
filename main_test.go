package main

import (
	"os"
	"testing"

	. "github.com/MrYZhou/outil/command"
	. "github.com/MrYZhou/outil/ssh"
)

func TestRun(t *testing.T) {
	Run(".", "docker stats")
}

func TestConnectWithKey(t *testing.T) {
	// 读取私钥文件内容
	contentBytes, _ := os.ReadFile("d:/larry.pem")
	var cli Cli
	cli.Host = "47.120.11.197:22"
	cli.User = "root"
	cli.PrivateKey = contentBytes
	con, _ := ConnectServer(cli)
	con.Run("cat /proc/cpuinfo")
}

func TestServer(t *testing.T) {
	var cli Cli
	cli.Host = "192.168.0.62:22"
	cli.User = "root"
	cli.Password = "YH4WfLbGPasRLVhs"
	con, _ := ConnectServer(cli)
	con.Run("cat /proc/cpuinfo")
}

func TestServerEasy(t *testing.T) {
	con, _ := Server("192.168.0.62:22","root","YH4WfLbGPasRLVhs")
	con.Run("cat /proc/cpuinfo")
}