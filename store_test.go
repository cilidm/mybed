package main

import (
	"mybedv2/app/system/model/store"
	storeService "mybedv2/app/system/service/store"
	"strings"
	"testing"
)

func TestStore(t *testing.T) {
	storeConf := store.GetStoreConfig("cloud_type", "cs-minio")
	cs, _ := storeService.NewCloudStoreByConf(storeConf, false)
	storePath := strings.ReplaceAll("http://127.0.0.1:9000/mybed/mybed/1/1032594_1.jpg", storeConf.PublicBucketDomain+storeConf.PublicBucket, "")
	t.Log(storePath)
	t.Log(cs.IsExist(storePath))
	if err := cs.Delete(storePath); err != nil {
		t.Error(err)
	}
}
