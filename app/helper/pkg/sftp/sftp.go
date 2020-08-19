package sftp

import (
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"mime"
	"mybedv2/app/helper/pkg/gsema"
	"mybedv2/app/helper/util/pathdir"
	img2 "mybedv2/app/system/model/img"
	"mybedv2/app/system/model/store"
	storeService "mybedv2/app/system/service/store"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
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
	User      string `json:"user" form:"user"`
	Pwd       string `json:"-" form:"pwd"`
	Host      string `json:"host" form:"host"`
	Port      int    `json:"port" form:"port"`
	Source    string `json:"source" form:"source"`
	KeepLocal int    `json:"keep_local"`
}

type Config struct {
	User string `json:"user"`
	Pwd  string `json:"pwd"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

type SftpConf struct {
	KeepLocal int    `json:"keep_local"`
	Uid       string `json:"uid"`
	ClientIP  string `json:"client_ip"`
	Source    string `json:"source"`
	Target    string `json:"target"`
	Config    Config
}

var (
	sema = gsema.NewSemaphore(20)
	cs   *store.CloudStore
)

func NewSftp(conf SftpConf) error {
	var (
		err        error
		sftpClient *sftp.Client
	)
	cs, err = storeService.NewCloudStore(false)
	if err != nil {
		return errors.New("未检测到可用的store")
	}
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
			sema.Add(1)
			go getFile(sftpClient, v, filepath.Join(conf.Target, stat.Name()), conf.Uid, conf.ClientIP, conf.KeepLocal)
		}
	}
	sema.Wait()
	return nil
}

func getFile(sftpClient *sftp.Client, remotePath, localPath, uid, ip string, keep int) {
	defer sema.Done()
	srcFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return
	}
	defer srcFile.Close()
	fmt.Println(srcFile.Name(), localPath)
	dstFile, err := os.Create(localPath)
	if err != nil {
		return
	}
	defer dstFile.Close()
	if _, err = srcFile.WriteTo(dstFile); err != nil {
		return
	}
	if err := StoreAndSave(localPath, uid, ip, keep); err != nil {
		return
	}
	return
}

func StoreAndSave(localPath, uid, ip string, keep int) error {
	md5, _ := pathdir.GetMd5V1(localPath)
	hasMd := img2.GetImgdataByMd5(md5)
	if hasMd.Id > 0 {
		return nil
	}

	stat, _ := os.Stat(localPath)
	miMe := mime.TypeByExtension(path.Ext(localPath))
	if err := cs.Upload(localPath, "sftp/"+uid+"/"+stat.Name(), map[string]string{"Content-Type": miMe}); err != nil {
		fmt.Println("上传到store出错", err.Error())
		return err
	}
	if keep == 2 {
		defer os.Remove(localPath)
	}
	var storeEntity store.Entity
	csc := storeEntity.FindOne()
	var mcPath string
	if csc.CloudType == "cs-minio" { //存储源minio路径需加bucket
		mcPath = csc.PublicBucketDomain + path.Join(csc.PublicBucket, "sftp/"+uid+"/"+stat.Name())
	} else {
		mcPath = csc.PublicBucketDomain + "sftp/" + uid + "/" + stat.Name()
	}
	uidInt, _ := strconv.ParseInt(uid, 10, 64)
	var img img2.Entity
	_, err := img.Insert(img2.Entity{
		ImgName:   stat.Name(),
		ImgUrl:    mcPath,
		UserId:    uidInt,
		Sizes:     stat.Size(),
		Source:    store.GetStoreNum(csc.CloudType),
		Md5:       md5,
		Abnormal:  ip,
		ImgType:   9,
		CreatedAt: time.Now(),
	})
	if err != nil {
		fmt.Println("保存img数据出错:", err.Error())
		return err
	}
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
