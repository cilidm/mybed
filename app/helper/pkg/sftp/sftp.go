package sftp

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"mybedv2/app/helper/util/pathdir"
	"net"
	"os"
	"path"
	"path/filepath"
	"time"
)

func connect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

type SftpBindForm struct {
	User   string `json:"user" form:"user"`
	Pwd    string `json:"pwd" form:"pwd"`
	Host   string `json:"host" form:"host"`
	Port   int    `json:"port" form:"port"`
	Source string `json:"source" form:"source"`
}

type Config struct {
	User string `json:"user"`
	Pwd  string `json:"pwd"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

type SftpConf struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Config Config
}

func NewSftp(conf SftpConf) error {
	var (
		err        error
		sftpClient *sftp.Client
	)
	pathdir.PathExists(conf.Target)
	// 这里换成实际的 SSH 连接的 用户名，密码，主机名或IP，SSH端口
	sftpClient, err = connect(conf.Config.User, conf.Config.Pwd, conf.Config.Host, conf.Config.Port)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	var dirName []string
	dirPath := conf.Source
	files, err := ShowDir(dirPath, dirName, sftpClient)
	for _, v := range files {
		stat, _ := sftpClient.Stat(v)
		if stat.IsDir() == false {
			getFile(sftpClient, v, filepath.Join(conf.Target, stat.Name()))
		}
	}
	return nil
}

func ShowDir(path string, dirName []string, client *sftp.Client) ([]string, error) {
	dir, err := client.Stat(path)
	if err != nil {
		return dirName, err
	}
	if dir.IsDir() {
		rd, err := client.ReadDir(path)
		if err != nil {
			return dirName, err
		}
		for _, v := range rd {
			if v.IsDir() {
				dirName, err = ShowDir(path+"/"+v.Name(), dirName, client)
				if err != nil {
					return dirName, err
				}
			} else {
				dirName = append(dirName, path+"/"+v.Name())
			}
		}
	}
	return dirName, nil
}

func getFile(sftpClient *sftp.Client, remotePath, localPath string) error {
	srcFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	fmt.Println(srcFile.Name(), localPath)
	dstFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err = srcFile.WriteTo(dstFile); err != nil {
		return err
	}
	//fmt.Println("copy " + srcFile.Name() + " from remote server finished!")
	return nil
}

func uploadFile(sftpClient *sftp.Client, remoteDir, localFilePath string) error {
	// 用来测试的本地文件路径 和 远程机器上的文件夹
	//var localFilePath = "/path/to/local/file/test.txt"
	//var remoteDir = "/remote/dir/"
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	var remoteFileName = srcFile.Name()
	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 5000)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}

	fmt.Println("copy file to remote server finished!")
	return nil
}
