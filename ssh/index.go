package ssh

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	. "github.com/MrYZhou/outil/common"
	. "github.com/MrYZhou/outil/file"
	"github.com/gosuri/uiprogress"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// 连接信息
type Cli struct {
	Host       string       // 主机地址
	User       string       // 登录用户
	Password   string       // 密码
	Client     *ssh.Client  // bash操作
	SftpClient *sftp.Client // 文件操作
	LastResult string       // 执行的最后一次结果
	PrivateKey []byte       // 私钥串
	PublicKey  []byte       // 远程服务器的公钥串,如果有传则校验,可以不传
}

// 连接对象
func (c *Cli) Connect() (*Cli, error) {
	config := &ssh.ClientConfig{}
	config.SetDefaults()
	config.User = c.User

	// 如果有私钥，按私钥连接,否则按密码
	if len(c.PrivateKey) > 0 {
		// 解析私钥
		signer, err := ssh.ParsePrivateKey(c.PrivateKey)
		if err != nil {
			log.Fatal("Failed to parse private key: ", err)
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		config.Auth = []ssh.AuthMethod{ssh.Password(c.Password)}
	}

	// 服务器公钥校验
	if len(c.PublicKey) > 0 {
		// 解析服务器公钥
		signer, err := ssh.ParsePublicKey(c.PublicKey)
		if err != nil {
			log.Fatal("Failed to parse private key: ", err)
		}
		config.HostKeyCallback = ssh.FixedHostKey(signer)
	} else {
		// 没有给公钥,使用ssh默认的一个不校验实现
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	// 客户端连接SSH服务器
	client, err := ssh.Dial("tcp", c.Host, config)
	sftp, err := sftp.NewClient(client)
	if nil != err {
		return c, err
	}
	c.Client = client
	c.SftpClient = sftp
	return c, nil
}

/*
获取服务器操作对象,简单版本。

如果需要通过密钥连接,请使用 ConnectServer
*/
func Server(host string, user string, password string) (*Cli, error) {

	cli := Cli{
		Host:     host,
		User:     user,
		Password: password,
	}
	c, err := cli.Connect()
	return c, err
}

// 获取服务器操作对象
func ConnectServer(cli Cli) (*Cli, error) {
	c, err := cli.Connect()
	return c, err
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

// 不需要输出信息
func (c Cli) RunQuiet(shell string) error {
	if c.Client == nil {
		if _, err := c.Connect(); err != nil {
			return err
		}
	}
	session, err := c.Client.NewSession()
	if err != nil {
		return err
	}
	// 关闭会话
	defer session.Close()
	return nil
}

/*
切片本地文件上传到远程

target 服务器的目录

filePath 切片的文件路径

num 切片数量
*/
func (c *Cli) SliceUpload(target string, filePath string, num int) []string {
	if num < 2 {
		fmt.Println("切片数量至少为2")
		return nil
	}
	c.createDir(target)

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("文件不存在")
		return nil
	}
	fileInfo, _ := f.Stat()

	defer f.Close()

	size := fileInfo.Size() / int64(num)
	duo := fileInfo.Size() - size*int64(num)
	fileList := make([]string, 0)

	var offset int64
	var offsetList []int64
	var chunkSizeList []int64

	for i := 0; i < num; i++ {
		if i == num-1 {
			size = size + duo
		}
		// 记录offset,和每一块的文件大小
		offsetList = append(offsetList, offset)
		chunkSizeList = append(chunkSizeList, size)
		offset += size

		rand_str := RandStr(10)
		targetPath := path.Join(target, "chunk"+rand_str)
		fileList = append(fileList, targetPath)
	}
	// 批量写入
	var wg sync.WaitGroup
	for i, targetPath := range fileList {
		wg.Add(1)
		go func(i int, targetPath string, f *os.File) {
			ftpFile, _ := c.SftpClient.Create(targetPath)
			size := chunkSizeList[i]
			offset := offsetList[i]
			chunk := make([]byte, size)

			f.ReadAt(chunk, offset)
			ftpFile.Write([]byte(chunk))

			wg.Done()
		}(i, targetPath, f)
	}
	wg.Wait()

	return fileList
}

/*
切片本地文件上传到远程,同时输出进度

target 服务器的目录

filePath 切片的文件路径

num 切片数量
*/
func (c *Cli) SliceUploadWithProgress(target string, filePath string, num int) []string {
	if num < 2 {
		fmt.Println("切片数量至少为2")
		return nil
	}
	c.createDir(target)

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("文件不存在")
		return nil
	}
	defer f.Close()

	fileInfo, _ := f.Stat()
	size := fileInfo.Size()
	duo := size - int64(num)*size/int64(num)
	chunkSize := size / int64(num)
	if duo > 0 {
		chunkSize += duo
	}

	fileList := make([]string, 0)
	offsetList := make([]int64, 0)
	chunkSizeList := make([]int64, num)

	for i := 0; i < num; i++ {
		if i == num-1 && duo > 0 {
			chunkSize = chunkSize + duo
		}
		offsetList = append(offsetList, int64(i) * chunkSize)
		chunkSizeList = append(chunkSizeList, chunkSize)
		rand_str := RandStr(10)
		targetPath := path.Join(target, "chunk"+rand_str)
		fileList = append(fileList, targetPath)
	}

	var wg sync.WaitGroup
	progress := uiprogress.New()
	bar := progress.AddBar(int(size)).AppendCompleted().PrependElapsed()
	progress.Start()

	uploadSizesCh := make(chan int64) // 创建一个用于传输已上传大小的通道

	for i, targetPath := range fileList {
		wg.Add(1)
		go func(i int, targetPath string, f *os.File, chunkSize int64) {
			ftpFile, _ := c.SftpClient.Create(targetPath)

			chunk := make([]byte, chunkSize)
			n, err := f.ReadAt(chunk, offsetList[i])
			if err != nil && err != io.EOF {
				fmt.Println("读取文件错误:", err)
				return
			}

			_, writeErr := ftpFile.Write(chunk[:n])
			if writeErr != nil {
				fmt.Println("上传文件片段错误:", writeErr)
				return
			}

			// 发送已上传的字节数量到通道
			uploadSizesCh <- int64(n)

			wg.Done()
		}(i, targetPath, f, chunkSizeList[i])
	}

	go func() {
		totalUploaded := int64(0)
		for range fileList {
			n := <-uploadSizesCh
			totalUploaded += n
			bar.Set(int(totalUploaded))
		}
	}()

	wg.Wait()

	close(uploadSizesCh)
	progress.Stop()

	return fileList
}

/*
关闭文件
*/
func (c *Cli) Close() {
	defer c.Client.Close()
	defer c.SftpClient.Close()
}

/*
合并远程文件
*/
func (c *Cli) ConcatRemoteFile(fileList []string, target string) {
	command := "cat "
	command += strings.Join(fileList, " ")
	command += " > "
	command += target
	c.Run(command)
}

/*
合并远程文件

fileList 文件列表

target 文件合成路径
*/
func (c *Cli) CombineRemoteFile(fileList []string, target string) {

	ftpFile, _ := c.CreateFile(target)
	defer ftpFile.Close()
	chunkList := make([][]byte, 0)
	sizeList := make([]int64, 0)
	var offset int64
	for _, name := range fileList {
		ftpBase, _ := c.SftpClient.Open(name)
		defer ftpBase.Close()

		fileInfo, _ := ftpBase.Stat()
		size := fileInfo.Size()
		buffer := make([]byte, size)
		ftpBase.Read(buffer)

		chunkList = append(chunkList, buffer)
		sizeList = append(sizeList, offset)
		offset += size

	}
	// 初始化文件大小
	alloc := make([]byte, offset)
	ftpFile.Write(alloc)
	// 并行写入
	var wg sync.WaitGroup
	for i, chunk := range chunkList {
		go wg.Add(1)
		func(i int) {
			offset := sizeList[i]
			ftpFile.WriteAt(chunk, offset)
			wg.Done()
		}(i)

	}
	wg.Wait()

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
	var wg sync.WaitGroup
	for i, f := range list {
		targetPath := strings.Replace(f, base, target, 1)
		wg.Add(1)
		go func(i int, f string, targetPath string) {
			c.UploadFile(f, targetPath)
			wg.Done()
		}(i, f, targetPath)
	}
	wg.Wait()
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
	ftpFile.Write(fileByte)

	defer ftpFile.Close()
	defer file.Close()
}
