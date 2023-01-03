package ssh

import (
	"io/ioutil"
	"net"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// 连接信息
type Cli struct {
	user       string
	password   string
	host       string
	client     *ssh.Client
	sftpClient *sftp.Client
	LastResult string
}

// 连接对象
func (c *Cli) Connect() (*Cli, error) {
	config := &ssh.ClientConfig{}
	config.SetDefaults()
	config.User = c.user
	config.Auth = []ssh.AuthMethod{ssh.Password(c.password)}
	config.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }
	client, err := ssh.Dial("tcp", c.host, config)
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

func Server(host string, user string, password string) Cli {

	cli := Cli{
		host:     host,
		user:     user,
		password: password,
	}
	c, _ := cli.Connect()

	defer c.client.Close()
	return cli
}
func (c Cli) UploadFile(localFile, remoteFileName string) {
	if c.client == nil {
		c.Connect();
	}
	file, _ := os.Open(localFile)
	defer file.Close()
	ftpFile, _ := c.sftpClient.Create(remoteFileName)
	defer ftpFile.Close()

	fileByte, _ := ioutil.ReadAll(ftpFile)
	ftpFile.Write(fileByte)
}
