package store

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"mybedv2/app/helper/e"
	"mybedv2/app/helper/pkg/cloudStore"
	"os"
	"strings"
)

type BindForm struct {
	AccessKey          string `json:"access_key" form:"access_key" binding:"required"`
	SecretKey          string `json:"secret_key" form:"secret_key" binding:"required"`
	Endpoint           string `json:"endpoint" form:"endpoint"`
	Region             string `json:"region" form:"region"`
	AppId              string `json:"app_id" form:"app-id"`
	PublicBucket       string `json:"public_bucket" form:"public_bucket" binding:"required"`
	PublicBucketDomain string `json:"public_bucket_domain" form:"public_bucket_domain" binding:"required"`
	CloudType          string `json:"cloud_type" form:"config_type" binding:"required"` //注:此处名称及绑定名称不一致
	Status             int    `json:"status"`
}

type ListShow struct {
	Id        int    `json:"id"`
	ImgName   string `json:"img_name"`
	ImgUrl    string `json:"img_url"`
	User      string `json:"user"`
	Size      string `json:"size"`
	Ip        string `json:"ip"`
	CreatedAt string `json:"created_at"`
}

func GetStoreNameByCloudType(str string) string {
	var storeConfig = map[string]string{
		"cs-oss":   "阿里云OSS",
		"cs-minio": "Minio",
		"cs-cos":   "腾讯云COS",
		"cs-obs":   "华为云OBS",
		"cs-bos":   "百度云BOS",
		"cs-upyun": "又拍云",
		"cs-qiniu": "七牛云",
	}
	return storeConfig[str]
}

func GetStoreNum(str string) int {
	var storeConfig = map[string]int{
		"cs-oss":   1,
		"cs-minio": 2,
		"cs-cos":   3,
		"cs-obs":   4,
		"cs-bos":   5,
		"cs-upyun": 6,
		"cs-qiniu": 7,
	}
	return storeConfig[str]
}

func GetStoreStr(num int) string {
	var storeConfig = map[int]string{
		1: "cs-oss",
		2: "cs-minio",
		3: "cs-cos",
		4: "cs-obs",
		5: "cs-bos",
		6: "cs-upyun",
		7: "cs-qiniu",
	}
	return storeConfig[num]
}

type ConfigOss struct {
	AccessKey           string `store:"access-key" fieldName:"AccessKey"`
	SecretKey           string `store:"secret-key" fieldName:"SecretKey"`
	Endpoint            string `store:"endpoint" fieldName:"Endpoint"`
	PublicBucket        string `store:"public-bucket" fieldName:"PublicBucket"`
	PublicBucketDomain  string `store:"public-bucket-domain" fieldName:"PublicBucketDomain"`
	PrivateBucket       string `store:"private-bucket" fieldName:"PrivateBucket"`
	PrivateBucketDomain string `store:"private-bucket-domain" fieldName:"PrivateBucketDomain"`
	Expire              int64  `store:"expire" fieldName:"Expire"`
}

type ConfigMinio struct {
	AccessKey           string `store:"access-key" fieldName:"AccessKey"`
	SecretKey           string `store:"secret-key" fieldName:"SecretKey"`
	Endpoint            string `store:"endpoint" fieldName:"Endpoint"`
	PublicBucket        string `store:"public-bucket" fieldName:"PublicBucket"`
	PublicBucketDomain  string `store:"public-bucket-domain" fieldName:"PublicBucketDomain"`
	PrivateBucket       string `store:"private-bucket" fieldName:"PrivateBucket"`
	PrivateBucketDomain string `store:"private-bucket-domain" fieldName:"PrivateBucketDomain"`
	Expire              int64  `store:"expire" fieldName:"Expire"`
}

type ConfigCos struct {
	AccessKey           string `store:"access-key" fieldName:"AccessKey"`
	SecretKey           string `store:"secret-key" fieldName:"SecretKey"`
	Region              string `store:"region" fieldName:"Region"`
	AppId               string `store:"app-id" fieldName:"AppId"`
	PublicBucket        string `store:"public-bucket" fieldName:"PublicBucket"`
	PublicBucketDomain  string `store:"public-bucket-domain" fieldName:"PublicBucketDomain"`
	PrivateBucket       string `store:"private-bucket" fieldName:"PrivateBucket"`
	PrivateBucketDomain string `store:"private-bucket-domain" fieldName:"PrivateBucketDomain"`
	Expire              int64  `store:"expire" fieldName:"Expire"`
}

type ConfigBos struct {
	AccessKey           string `store:"access-key" fieldName:"AccessKey"`
	SecretKey           string `store:"secret-key" fieldName:"SecretKey"`
	Endpoint            string `store:"endpoint" fieldName:"Endpoint"`
	PublicBucket        string `store:"public-bucket" fieldName:"PublicBucket"`
	PublicBucketDomain  string `store:"public-bucket-domain" fieldName:"PublicBucketDomain"`
	PrivateBucket       string `store:"private-bucket" fieldName:"PrivateBucket"`
	PrivateBucketDomain string `store:"private-bucket-domain" fieldName:"PrivateBucketDomain"`
	Expire              int64  `store:"expire" fieldName:"Expire"`
}

type ConfigObs struct {
	AccessKey           string `store:"access-key" fieldName:"AccessKey"`
	SecretKey           string `store:"secret-key" fieldName:"SecretKey"`
	Endpoint            string `store:"endpoint" fieldName:"Endpoint"`
	PublicBucket        string `store:"public-bucket" fieldName:"PublicBucket"`
	PublicBucketDomain  string `store:"public-bucket-domain" fieldName:"PublicBucketDomain"`
	PrivateBucket       string `store:"private-bucket" fieldName:"PrivateBucket"`
	PrivateBucketDomain string `store:"private-bucket-domain" fieldName:"PrivateBucketDomain"`
	Expire              int64  `store:"expire" fieldName:"Expire"`
}

type ConfigQiniu struct {
	AccessKey           string `store:"access-key" fieldName:"AccessKey"`
	SecretKey           string `store:"secret-key" fieldName:"SecretKey"`
	Endpoint            string `store:"endpoint" fieldName:"Endpoint"`
	PublicBucket        string `store:"public-bucket" fieldName:"PublicBucket"`
	PublicBucketDomain  string `store:"public-bucket-domain" fieldName:"PublicBucketDomain"`
	PrivateBucket       string `store:"private-bucket" fieldName:"PrivateBucket"`
	PrivateBucketDomain string `store:"private-bucket-domain" fieldName:"PrivateBucketDomain"`
	Expire              int64  `store:"expire" fieldName:"Expire"`
}

type ConfigUpYun struct {
	AccessKey           string `store:"access-key" fieldName:"AccessKey"`
	SecretKey           string `store:"secret-key" fieldName:"SecretKey"`
	Endpoint            string `store:"endpoint" fieldName:"Endpoint"`
	PublicBucket        string `store:"public-bucket" fieldName:"PublicBucket"`
	PublicBucketDomain  string `store:"public-bucket-domain" fieldName:"PublicBucketDomain"`
	PrivateBucket       string `store:"private-bucket" fieldName:"PrivateBucket"`
	PrivateBucketDomain string `store:"private-bucket-domain" fieldName:"PrivateBucketDomain"`
	Expire              int64  `store:"expire" fieldName:"Expire"`
}

type CloudStore struct {
	Private       bool
	StoreType     string
	CanGZIP       bool
	Client        interface{}
	Config        interface{}
	Expire        int64
	PublicDomain  string
	PrivateDomain string
}

func (c *CloudStore) Lists(object string) (files []cloudStore.File, err error) {
	switch c.StoreType {
	case e.StoreCos:
		files, err = c.Client.(*cloudStore.COS).Lists(object)
	case e.StoreOss:
		files, err = c.Client.(*cloudStore.OSS).Lists(object)
	case e.StoreBos:
		files, err = c.Client.(*cloudStore.BOS).Lists(object)
	case e.StoreObs:
		files, err = c.Client.(*cloudStore.OBS).Lists(object)
	case e.StoreUpyun:
		files, err = c.Client.(*cloudStore.UpYun).Lists(object)
	case e.StoreMinio:
		files, err = c.Client.(*cloudStore.MinIO).Lists(object)
	case e.StoreQiniu:
		files, err = c.Client.(*cloudStore.QINIU).Lists(object)
	}
	return
}

func (c *CloudStore) Upload(tmpFile, saveFile string, headers ...map[string]string) (err error) {
	switch c.StoreType {
	case e.StoreCos:
		err = c.Client.(*cloudStore.COS).Upload(tmpFile, saveFile, headers...)
	case e.StoreOss:
		err = c.Client.(*cloudStore.OSS).Upload(tmpFile, saveFile, headers...)
	case e.StoreBos:
		err = c.Client.(*cloudStore.BOS).Upload(tmpFile, saveFile, headers...)
	case e.StoreObs:
		err = c.Client.(*cloudStore.OBS).Upload(tmpFile, saveFile, headers...)
	case e.StoreUpyun:
		err = c.Client.(*cloudStore.UpYun).Upload(tmpFile, saveFile, headers...)
	case e.StoreMinio:
		err = c.Client.(*cloudStore.MinIO).Upload(tmpFile, saveFile, headers...)
	case e.StoreQiniu:
		err = c.Client.(*cloudStore.QINIU).Upload(tmpFile, saveFile, headers...)
	}
	return
}

func (c *CloudStore) Delete(objects ...string) (err error) {
	switch c.StoreType {
	case e.StoreCos:
		err = c.Client.(*cloudStore.COS).Delete(objects...)
	case e.StoreOss:
		err = c.Client.(*cloudStore.OSS).Delete(objects...)
	case e.StoreBos:
		err = c.Client.(*cloudStore.BOS).Delete(objects...)
	case e.StoreObs:
		err = c.Client.(*cloudStore.OBS).Delete(objects...)
	case e.StoreUpyun:
		err = c.Client.(*cloudStore.UpYun).Delete(objects...)
	case e.StoreMinio:
		err = c.Client.(*cloudStore.MinIO).Delete(objects...)
	case e.StoreQiniu:
		err = c.Client.(*cloudStore.QINIU).Delete(objects...)
	}
	return
}

// err 返回 nil，表示文件存在，否则表示文件不存在
func (c *CloudStore) IsExist(object string) (err error) {
	switch c.StoreType {
	case e.StoreCos:
		err = c.Client.(*cloudStore.COS).IsExist(object)
	case e.StoreOss:
		err = c.Client.(*cloudStore.OSS).IsExist(object)
	case e.StoreBos:
		err = c.Client.(*cloudStore.BOS).IsExist(object)
	case e.StoreObs:
		err = c.Client.(*cloudStore.OBS).IsExist(object)
	case e.StoreUpyun:
		err = c.Client.(*cloudStore.UpYun).IsExist(object)
	case e.StoreMinio:
		err = c.Client.(*cloudStore.MinIO).IsExist(object)
	case e.StoreQiniu:
		err = c.Client.(*cloudStore.QINIU).IsExist(object)
	}
	return
}

func (c *CloudStore) PingTest() (err error) {
	tmpFile := "mybed-test-file.txt"
	saveFile := "mybed-test-file.txt"
	text := "hello world"

	defer func() {
		if err != nil {
			err = fmt.Errorf("Bucket是否私有：%v，错误信息：%v", c.Private, err.Error())
		}
	}()

	err = ioutil.WriteFile(tmpFile, []byte(text), os.ModePerm)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	if err = c.Upload(tmpFile, saveFile); err != nil {
		return
	}
	if err = c.IsExist(saveFile); err != nil {
		return
	}
	if err = c.Delete(saveFile); err != nil {
		return
	}
	return
}

//设置默认图片
//@param                picture             图片文件或者图片文件md5等
//@param                ext                 图片扩展名，如果图片文件参数(picture)的值为md5时，需要加上后缀扩展名
//@return               link                图片url链接
func (c *CloudStore) getImageFromCloudStore(picture string, ext ...string) (link string) {
	if len(ext) > 0 {
		picture = picture + "." + ext[0]
	} else if !strings.Contains(picture, ".") && len(picture) > 0 {
		picture = picture + ".jpg"
	}
	if c == nil || c.Client == nil {
		return
	}

	return c.GetSignURL(picture)
}

func (c *CloudStore) GetSignURL(object string) (link string) {
	var err error
	switch c.StoreType {
	case e.StoreCos:
		link, err = c.Client.(*cloudStore.COS).GetSignURL(object, c.Expire)
	case e.StoreOss:
		link, err = c.Client.(*cloudStore.OSS).GetSignURL(object, c.Expire)
	case e.StoreBos:
		link, err = c.Client.(*cloudStore.BOS).GetSignURL(object, c.Expire)
	case e.StoreObs:
		link, err = c.Client.(*cloudStore.OBS).GetSignURL(object, c.Expire)
	case e.StoreUpyun:
		link, err = c.Client.(*cloudStore.UpYun).GetSignURL(object, c.Expire)
	case e.StoreMinio:
		link, err = c.Client.(*cloudStore.MinIO).GetSignURL(object, c.Expire)
	case e.StoreQiniu:
		link, err = c.Client.(*cloudStore.QINIU).GetSignURL(object, c.Expire)
	}
	if err != nil {
		fmt.Println(err)
	}
	return
}

func (c *CloudStore) ImageWithDomain(htmlOld string) (htmlNew string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlOld))
	if err != nil {
		fmt.Println(err)
		return htmlOld
	}
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		if src, exist := s.Attr("src"); exist {
			if !(strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://")) {
				s.SetAttr("src", c.GetSignURL(src))
			}
		}

	})
	htmlNew, _ = doc.Find("body").Html()
	return
}

func (c *CloudStore) ImageWithoutDomain(htmlOld string) (htmlNew string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlOld))
	if err != nil {
		fmt.Println(err)
		return htmlOld
	}
	domain := c.GetPublicDomain()

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		if src, exist := s.Attr("src"); exist {
			//不存在http开头的图片链接，则更新为绝对链接
			if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
				src = strings.TrimPrefix(src, domain)
				s.SetAttr("src", src)
			}
		}
	})
	htmlNew, _ = doc.Find("body").Html()
	return
}

//从HTML中提取图片文件，并删除
func (c *CloudStore) DeleteImageFromHtml(htmlStr string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		fmt.Println(err)
		return
	}

	var objects []string

	domain := c.GetPublicDomain()

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		if src, exist := s.Attr("src"); exist {
			//不存在http开头的图片链接，则更新为绝对链接
			if !(strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://")) {
				objects = append(objects, src)
			} else {
				src = strings.TrimPrefix(src, domain)
				objects = append(objects, src)
			}
		}
	})
	if err = c.Delete(objects...); err != nil {
		fmt.Println(err)
	}
}

func (c *CloudStore) GetPublicDomain() (domain string) {
	object := "test.mybed.test"
	link := c.GetSignURL(object)
	return strings.TrimRight(strings.Split(link, object)[0], "/")
}

type StoreMenu struct {
	Name string
	Path string
}

func GetStoreMenu() []StoreMenu {
	var menus []StoreMenu
	storeMenu := GetNotNullConfig()
	for _, v := range storeMenu {
		menus = append(menus, StoreMenu{
			Name: GetStoreNameByCloudType(v.CloudType),
			Path: v.CloudType,
		})
	}
	return menus
}
