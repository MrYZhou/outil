package ssh

import (
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

// 连接信息
type Cli struct {
	user       string
	password   string
	addr       string
	client     *ssh.Client
	LastResult string
}

// 连接对象
func (c *Cli) Connect() (*Cli, error) {
	config := &ssh.ClientConfig{}
	config.SetDefaults()
	config.User = c.user
	config.Auth = []ssh.AuthMethod{ssh.Password(c.password)}
	config.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }
	client, err := ssh.Dial("tcp", c.addr, config)
	if nil != err {
		return c, err
	}
	c.client = client
	return c, nil
}

// 执行shell
func (c Cli) Run(shell string) (string, error) {
	if c.client == nil {
		if _, err := c.Connect(); err != nil {
			return "", err
		}
	}
	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	// 关闭会话
	defer session.Close()
	buf, err := session.CombinedOutput(shell)

	c.LastResult = string(buf)
	return c.LastResult, err
}

func Server() Cli {
	addr := os.Getenv("addr")
	user := os.Getenv("user")
	password := os.Getenv("password")
	cli := Cli{
		addr: addr,
		user: user,
		password:password,
	}
	c, _ := cli.Connect()

	defer c.client.Close()
	return cli
}

// func Server(addr string,user string, password string) Cli {

// 	cli := Cli{
// 		addr: addr,
// 		user: user,
// 		password:password,
// 	}
// 	c, _ := cli.Connect()

// 	defer c.client.Close()
// 	return cli
// }
