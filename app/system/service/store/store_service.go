package store

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mime"
	"mybedv2/app/helper/e"
	"mybedv2/app/helper/pkg/cloudStore"
	storeModel "mybedv2/app/system/model/store"
	"mybedv2/app/system/service/user"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
)

type ConfigCate string

// 上传图片到默认存储源 filePath 图片地址 存储路径：StroeDefaultDir(这里是mybed) / userID / 图片名称
func DefaultUploadStore(filePath string, c *gin.Context) (string, error) {
	cs, err := NewCloudStore(false)
	if err != nil {
		return "", err
	}
	s, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}
	storePath := filepath.Join(e.StroeDefaultDir, strconv.FormatInt(user.GetUserId(c), 10)+"/"+s.Name())
	miMe := mime.TypeByExtension(path.Ext(filePath))
	err = cs.Upload(filePath, storePath, map[string]string{"Content-Type": miMe})
	return storePath, err
}

//传入CloudStoreConfig结构体返回通用CloudStore
func NewCloudStoreByConf(csc storeModel.Entity, private bool) (cs *storeModel.CloudStore, err error) {
	storeType := csc.CloudType
	config := GetConfigType(csc)
	private = false
	return NewCloudStoreWithConfig(config, storeType, private)
}

func GetConfigType(config storeModel.Entity) (cfg interface{}) {
	switch config.CloudType {
	case e.StoreCos:
		cfg = &storeModel.ConfigCos{}
	case e.StoreBos:
		cfg = &storeModel.ConfigBos{}
	case e.StoreOss:
		cfg = &storeModel.ConfigOss{}
	case e.StoreMinio:
		cfg = &storeModel.ConfigMinio{}
	case e.StoreUpyun:
		cfg = &storeModel.ConfigUpYun{}
	case e.StoreQiniu:
		cfg = &storeModel.ConfigQiniu{}
	case e.StoreObs:
		cfg = &storeModel.ConfigObs{}
	}
	vc := reflect.ValueOf(config)
	t := reflect.TypeOf(cfg)
	v := reflect.ValueOf(cfg)
	numFields := t.Elem().NumField()
	for i := 0; i < numFields; i++ {
		fieldName := t.Elem().Field(i).Tag.Get("fieldName")
		val := vc.FieldByName(fieldName)
		v.Elem().FieldByName(fieldName).Set(val)
	}
	return
}

// 创建云存储
//config 相应存储的struct指针
//storetype 云存储对应的名字
//private false不使用普通网络
func NewCloudStore(private bool) (cs *storeModel.CloudStore, err error) {
	var storeEntity storeModel.Entity
	csc := storeEntity.FindOne()
	storeType := csc.CloudType
	config := GetConfigType(csc)
	private = false
	return NewCloudStoreWithConfig(config, storeType, private)
}

var errWithoutConfig = errors.New("云存储配置不正确")

func NewCloudStoreWithConfig(storeConfig interface{}, storeType string, private bool) (cs *storeModel.CloudStore, err error) {
	cs = &storeModel.CloudStore{
		StoreType: storeType,
		Config:    storeConfig,
	}
	cs.Private = private
	switch cs.StoreType {
	case e.StoreOss:
		cfg := cs.Config.(*storeModel.ConfigOss)
		bucket := cfg.PublicBucket
		domain := cfg.PublicBucketDomain
		if cs.Private {
			bucket = cfg.PrivateBucket
			domain = cfg.PrivateBucketDomain
			if cfg.Expire <= 0 {
				cfg.Expire = 1800
			}
			cs.Expire = cfg.Expire
		}
		cs.PrivateDomain = cfg.PrivateBucketDomain
		cs.PublicDomain = cfg.PublicBucketDomain
		if cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Endpoint == "" || bucket == "" {
			err = errWithoutConfig
			return
		}
		cs.Client, err = cloudStore.NewOSS(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint, bucket, domain)
		cs.CanGZIP = true
	case e.StoreObs:
		cfg := cs.Config.(*storeModel.ConfigObs)
		bucket := cfg.PublicBucket
		domain := cfg.PublicBucketDomain
		if cs.Private {
			bucket = cfg.PrivateBucket
			domain = cfg.PrivateBucketDomain
			if cfg.Expire <= 0 {
				cfg.Expire = 1800
			}
			cs.Expire = cfg.Expire
		}
		cs.PrivateDomain = cfg.PrivateBucketDomain
		cs.PublicDomain = cfg.PublicBucketDomain
		if cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Endpoint == "" || bucket == "" {
			err = errWithoutConfig
			return
		}
		cs.Client, err = cloudStore.NewOBS(cfg.AccessKey, cfg.SecretKey, bucket, cfg.Endpoint, domain)
	case e.StoreQiniu:
		cfg := cs.Config.(*storeModel.ConfigQiniu)
		bucket := cfg.PublicBucket
		domain := cfg.PublicBucketDomain
		if cs.Private {
			bucket = cfg.PrivateBucket
			domain = cfg.PrivateBucketDomain
			if cfg.Expire <= 0 {
				cfg.Expire = 1800
			}
			cs.Expire = cfg.Expire
		}
		cs.PrivateDomain = cfg.PrivateBucketDomain
		cs.PublicDomain = cfg.PublicBucketDomain
		if cfg.AccessKey == "" || cfg.SecretKey == "" || bucket == "" {
			err = errWithoutConfig
			return
		}
		cs.Client, err = cloudStore.NewQINIU(cfg.AccessKey, cfg.SecretKey, bucket, domain)
	case e.StoreUpyun:
		cfg := cs.Config.(*storeModel.ConfigUpYun)
		bucket := cfg.PublicBucket
		domain := cfg.PublicBucketDomain
		if cs.Private {
			bucket = cfg.PrivateBucket
			domain = cfg.PrivateBucketDomain
			if cfg.Expire <= 0 {
				cfg.Expire = 1800
			}
			cs.Expire = cfg.Expire
		}
		cs.PrivateDomain = cfg.PrivateBucketDomain
		cs.PublicDomain = cfg.PublicBucketDomain
		//if cfg.Operator == "" || cfg.Password == "" || bucket == "" {
		if cfg.AccessKey == "" || cfg.SecretKey == "" || bucket == "" {
			err = errWithoutConfig
			return
		}
		//cs.client = CloudStore2.NewUpYun(bucket, cfg.Operator, cfg.Password, domain, cfg.Secret)
		cs.Client = cloudStore.NewUpYun(bucket, cfg.AccessKey, cfg.SecretKey, domain, cfg.Endpoint)
	case e.StoreMinio:
		cfg := cs.Config.(*storeModel.ConfigMinio)
		bucket := cfg.PublicBucket
		domain := cfg.PublicBucketDomain
		if cs.Private {
			bucket = cfg.PrivateBucket
			domain = cfg.PrivateBucketDomain
			if cfg.Expire <= 0 {
				cfg.Expire = 1800
			}
			cs.Expire = cfg.Expire
		}
		cs.PrivateDomain = cfg.PrivateBucketDomain
		cs.PublicDomain = cfg.PublicBucketDomain
		if cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Endpoint == "" || bucket == "" {
			err = errWithoutConfig
			return
		}
		cs.Client, err = cloudStore.NewMinIO(cfg.AccessKey, cfg.SecretKey, bucket, cfg.Endpoint, domain)
		cs.CanGZIP = true
	case e.StoreBos:
		cfg := cs.Config.(*storeModel.ConfigBos)
		bucket := cfg.PublicBucket
		domain := cfg.PublicBucketDomain
		if cs.Private {
			bucket = cfg.PrivateBucket
			domain = cfg.PrivateBucketDomain
			if cfg.Expire <= 0 {
				cfg.Expire = 1800
			}
			cs.Expire = cfg.Expire
		}
		cs.PrivateDomain = cfg.PrivateBucketDomain
		cs.PublicDomain = cfg.PublicBucketDomain
		if cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Endpoint == "" || bucket == "" {
			err = errWithoutConfig
			return
		}
		cs.Client, err = cloudStore.NewBOS(cfg.AccessKey, cfg.SecretKey, bucket, cfg.Endpoint, domain)
		cs.CanGZIP = true
	case e.StoreCos:
		cfg := cs.Config.(*storeModel.ConfigCos)
		bucket := cfg.PublicBucket
		domain := cfg.PublicBucketDomain
		if cs.Private {
			bucket = cfg.PrivateBucket
			domain = cfg.PrivateBucketDomain
			if cfg.Expire <= 0 {
				cfg.Expire = 1800
			}
			cs.Expire = cfg.Expire
		}
		cs.PrivateDomain = cfg.PrivateBucketDomain
		cs.PublicDomain = cfg.PublicBucketDomain
		if cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.AppId == "" || bucket == "" || cfg.Region == "" {
			err = errWithoutConfig
			return
		}
		cs.Client, err = cloudStore.NewCOS(cfg.AccessKey, cfg.SecretKey, bucket, cfg.AppId, cfg.Region, domain)
		cs.CanGZIP = true
	}
	return
}
