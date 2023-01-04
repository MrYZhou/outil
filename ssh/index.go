package ssh

import (
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strings"

	. "github.com/MrYZhou/outil/file"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// 连接信息
type Cli struct {
	user       string
	password   string
	host       string
	Client     *ssh.Client
	SftpClient *sftp.Client
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
	sftp, err := sftp.NewClient(client)
	if nil != err {
		return c, err
	}
	c.Client = client
	c.SftpClient = sftp
	return c, nil
}

// 执行shellclient
func (c Cli) Run(shell string) (string, error) {
	if c.Client == nil {
		if _, err := c.Connect(); err != nil {
			return "", err
		}
	}
	session, err := c.Client.NewSession()
	if err != nil {
		return "", err
	}
	// 关闭会话
	defer session.Close()

	r, err := session.StdoutPipe()
	if err != nil {
			fmt.Println(err)
			os.Exit(1001)
	}
	go io.Copy(os.Stdout, r)

	buf, err := session.CombinedOutput(shell)
	c.LastResult = string(buf)
	return c.LastResult, err
}

func Server(host string, user string, password string) *Cli {

	cli := Cli{
		host:     host,
		user:     user,
		password: password,
	}
	c, _ := cli.Connect()
	return c
}

// 创建目录
func (c *Cli) createDir(dir string) {
	c.SftpClient.MkdirAll(dir)
}

// 批量创建目录
func (c *Cli) createDirList(list []string) {
	for _, dir := range list {
		c.createDir(dir)
	}
}

// 判断文件是否存在
func (c *Cli) IsFileExist(path string) bool {
	info, _ := c.SftpClient.Stat(path)
	if info != nil {
		return true
	}
	return false
}

/*
创建文件

remoteFileName 文件名
*/
func (c *Cli) CreateFile(remoteFileName string) (*sftp.File, error) {
	remoteDir := path.Dir(remoteFileName)
	c.SftpClient.MkdirAll(remoteDir)
	ftpFile, err := c.SftpClient.Create(remoteFileName)
	return ftpFile, err
}
func initClient(c *Cli) *Cli {

	if c.SftpClient == nil {
		cli, _ := c.Connect()
		return cli
	} else {
		return c
	}

}

/*
上传目录到服务器

base 本地文件夹路径

target 远程文件夹路径
*/
func (c *Cli) UploadDir(base string, target string) {
	c.SftpClient.MkdirAll(target)
	list, dirList := ReadDirAll(base)
	// 创建远程目录
	for _, f := range dirList {
		targetPath := strings.Replace(f, base, target, 1)
		c.SftpClient.MkdirAll(targetPath)
	}
	// 创建远程文件
	for _, f := range list {
		targetPath := strings.Replace(f, base, target, 1)
		c.UploadFile(f, targetPath)
	}
}

/*
上传远程文件

localFile 本地文件路径

remoteFileName 远程文件路径
*/
func (c *Cli) UploadFile(localFile, remoteFileName string) {

	file, _ := os.Open(localFile)
	ftpFile, err := c.CreateFile(remoteFileName)
	if nil != err {
		fmt.Println(err)
	}

	fileByte, _ := io.ReadAll(file)
	go ftpFile.Write(fileByte)

	defer ftpFile.Close()
	defer file.Close()
}
