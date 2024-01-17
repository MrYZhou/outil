package main

import (
	"fmt"
	"testing"

	. "github.com/MrYZhou/outil/command"
	. "github.com/MrYZhou/outil/ssh"
)
func TestRun(t *testing.T) {
	Run(".","docker stats")
}

func TestConnectWithKey(t *testing.T) {
	var cli Cli
	cli.PrivateKey = "C:\\Users\\yzhou\\Desktop\\key\\id_rsa"
	con, err := ConnectServer(cli)
	fmt.Println(con, err)
}
