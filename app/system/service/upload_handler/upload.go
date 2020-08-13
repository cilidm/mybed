package upload_handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"mime"
	"mime/multipart"
	"mybedv2/app/helper/baidu"
	"mybedv2/app/helper/e"
	"mybedv2/app/helper/redis"
	"mybedv2/app/helper/util/pathdir"
	"mybedv2/app/helper/util/str"
	"mybedv2/app/system/model"
	"mybedv2/app/system/model/bd/bd_img"
	"mybedv2/app/system/model/blacklist/img_blacklist"
	imgModel "mybedv2/app/system/model/img"
	"mybedv2/app/system/model/store"
	imgService "mybedv2/app/system/service/img"
	storeService "mybedv2/app/system/service/store"
	"mybedv2/app/system/service/user"
	"mybedv2/conf"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

var LinkBucketStore = []string{e.StoreMinio}

func UploadHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, model.UploadResp{Code: strconv.Itoa(e.ErrorUploadForm), Info: e.GetMsg(e.ErrorUploadForm)})
		return
	}
	files := form.File["file"]
	resp, err := UploadImpl(files, c)
	if err != nil {
		if resp.Code == strconv.Itoa(e.NotEnoughFreeSpace) {
			c.JSON(http.StatusOK, model.UploadResp{Code: resp.Code, Info: resp.Info})
		} else {
			c.JSON(http.StatusBadRequest, model.UploadResp{Code: resp.Code, Info: resp.Info})
		}
		return
	}
	//csc := models.GetUseCloudStore()
	//SmallPngAndStoreImpl(csc) //缩略图
	c.JSON(http.StatusOK, model.UploadResp{Code: strconv.Itoa(e.SUCCESS), Data: resp.Data, Info: e.GetMsg(e.SUCCESS)})
}

// todo 清理废弃的文件夹 将临时文件夹路径及过期时间放入redis list，每天清理一次
func UploadImpl(files []*multipart.FileHeader, c *gin.Context) (model.UploadResp, error) {
	var (
		resp model.UploadResp
	)
	// 上传临时文件夹 static/upload/日期小时
	uploadTmp := path.Join(conf.Setting.UploadTmpDir, time.Now().Format(e.HourTimeFormatDir))
	pathdir.PathExists(uploadTmp)
	redis.Client.SAdd(e.TempDirList, uploadTmp) //加入待删除文件夹 每天删除一次
	// 获取所有图片
	for _, file := range files {
		//获取userid  如果是同一个人传的图片已存在 则直接返回 如果是不同人传的图片 则直接入库 user_imgdata
		uid := user.GetUserId(c)
		// 用户存储配额超限 游客上传不查配额 (uid = 0)
		if uid != 0 && imgService.IsMaxSize(uid, file.Size) == true { //webuploade/js/upload.js line201:
			return model.UploadResp{Code: strconv.Itoa(e.NotEnoughFreeSpace), Info: e.GetMsg(e.NotEnoughFreeSpace)}, errors.New("空间超限")
		}
		// 逐个存
		fileTmpPath := path.Join(uploadTmp, file.Filename) //原图的全路径
		if err := c.SaveUploadedFile(file, fileTmpPath); err != nil {
			return model.UploadResp{Code: strconv.Itoa(e.ErrorUploadSave), Info: e.GetMsg(e.ErrorUploadSave)}, err
		}
		md5, _ := pathdir.GetMd5V1(fileTmpPath)
		blackList := img_blacklist.FindByMd5(md5)
		if blackList > 0 { //在图片黑名单
			return model.UploadResp{Code: strconv.Itoa(e.ErrorUploadSave), Info: e.GetMsg(e.ErrorUploadSave)}, errors.New("此图片禁止上传")
		}
		imgExamine := bd_img.FindOne()
		if imgExamine.Status == 1 {
			bdResp, err := baidu.GetIdentificationResult(fileTmpPath, 1)
			if err != nil {
				return model.UploadResp{Code: strconv.Itoa(e.ErrorUploadSave), Info: e.GetMsg(e.ErrorUploadSave)}, err
			}
			if bdResp.ConclusionType != 1 {
				fs, _ := os.Stat(fileTmpPath)
				var blist img_blacklist.Entity
				err := blist.Insert(img_blacklist.BindForm{
					Info:     bdResp.ErrorMsg + bdResp.Conclusion,
					UserIp:   c.ClientIP(),
					UserId:   int(user.GetUserId(c)),
					FileName: fs.Name(),
					FileSize: fs.Size(),
					FileMd5:  md5,
				})
				if err != nil {
					return model.UploadResp{Code: strconv.Itoa(e.ErrorUploadSave), Info: e.GetMsg(e.ErrorUploadSave)}, errors.New(err.Error() + bdResp.ErrorMsg + bdResp.Conclusion)
				}
				return model.UploadResp{Code: strconv.Itoa(e.ErrorUploadSave), Info: e.GetMsg(e.ErrorUploadSave)}, errors.New(bdResp.ErrorMsg + bdResp.Conclusion)
			}
		}

		r, err := uploadSave(fileTmpPath, file.Filename, uid, c, &resp)
		defer os.Remove(fileTmpPath)
		if err != nil {
			return r, err
		}
	}
	return resp, nil
}

func uploadSave(fileTmpPath, fileName string, uid int64, c *gin.Context, resp *model.UploadResp) (model.UploadResp, error) {
	var imgInfo model.ImgInfo
	md5, _ := pathdir.GetMd5V1(fileTmpPath)
	img := imgModel.GetImgdataByMd5(md5)
	if img.Id > 0 { //已经存在同样的图片
		os.Remove(fileTmpPath)
		if img.UserId != uid {
			//图片信息入库
			img.UserId = uid
			img.Abnormal = c.ClientIP()
			_, err := imgModel.CreateImgdataByMd5(img)
			if err != nil {
				return model.UploadResp{Code: strconv.Itoa(e.ErrorSaveImgdata), Info: e.GetMsg(e.ErrorSaveImgdata) + err.Error()}, err
			}
		}
		//直接把数据库里图片名称及url放入数组里
		resp.Data = append(resp.Data, model.ImgInfo{Imgnames: img.ImgName, Imgurls: img.ImgUrl})
	} else {
		//上传原图
		// todo 将store初始化方法解藕
		cs, err := storeService.NewCloudStore(false)
		if err != nil {
			return model.UploadResp{Code: strconv.Itoa(e.ErrorUploadStore), Info: e.GetMsg(e.ErrorUploadStore) + err.Error()}, err
		}
		//store路径 userid + file.Filename
		storePath := filepath.Join(e.StroeDefaultDir, strconv.FormatInt(uid, 10)+"/"+fileName)
		//直接写入store
		miMe := mime.TypeByExtension(path.Ext(fileTmpPath))
		if err = cs.Upload(fileTmpPath, storePath, map[string]string{"Content-Type": miMe}); err != nil {
			return model.UploadResp{Code: strconv.Itoa(e.ErrorUploadStore), Info: e.GetMsg(e.ErrorUploadStore) + err.Error()}, err
		}
		var storeEntity store.Entity
		csc := storeEntity.FindOne()
		// store地址: ip/域名 + bucket + storepath
		var mcPath string
		if str.IsContain(LinkBucketStore, csc.CloudType) { //存储源minio路径需加bucket
			mcPath = csc.PublicBucketDomain + path.Join(csc.PublicBucket, storePath)
		} else {
			mcPath = csc.PublicBucketDomain + storePath
		}
		//图片信息入库
		//imgId, err := imgService.CreateImgdata(fileTmpPath, mcPath, c.ClientIP(), md5, uid, csc.CloudType)
		_, err = imgService.CreateImgdata(fileTmpPath, mcPath, c.ClientIP(), md5, uid, csc.CloudType)
		if err != nil {
			return model.UploadResp{Code: strconv.Itoa(e.ErrorSaveImgdata), Info: e.GetMsg(e.ErrorSaveImgdata) + err.Error()}, err
		}
		imgInfo.Imgnames = fileName
		imgInfo.Imgurls = mcPath
		resp.Data = append(resp.Data, imgInfo)
		//缩略图
	}
	return *resp, nil
}
