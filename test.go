package main

import (
	"fmt"
	"mybedv2/app/system/model/code"
	"mybedv2/app/system/model/store"
	storeService "mybedv2/app/system/service/store"
	"strings"
	"time"
)

func testStore() {
	storeConf := store.GetStoreConfig("cloud_type", "cs-minio")
	cs, _ := storeService.NewCloudStoreByConf(storeConf, false)
	storePath := strings.ReplaceAll("http://127.0.0.1:9000/mybed/mybed/1/1032594_1.jpg", storeConf.PublicBucketDomain+storeConf.PublicBucket, "")
	fmt.Println(storePath)
	fmt.Println(cs.IsExist(storePath))
	if err := cs.Delete(storePath); err != nil {
		fmt.Println(err)
	}
}

func testMany2Many() {
	var codeEntity code.Entity
	codeEntity.Insert(code.BindForm{
		Code:      "aaa",
		Value:     100,
		Status:    1,
		CreatedAt: time.Now(),
		UserId:    1,
	})
}

func main() {
	//testMany2Many()
}
