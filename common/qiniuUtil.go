package common

import (
	"context"
	"fmt"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/rs/xid"
	"mogu-go-v2/models"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/5 2:26 下午
 * @version 1.0
 */

type qiniuUtil struct{}

func (qiniuUtil) UploadQiNiu(localFilePath string, qiNiuConfig models.SystemConfig) string {
	accessKey := qiNiuConfig.QiNiuAccessKey
	secretKey := qiNiuConfig.QiNiuSecretKey
	bucket := qiNiuConfig.QiNiuBucket
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = setQiNiuArea(qiNiuConfig.QiNiuArea)
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = true
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	key := "mogu-go/" + xid.New().String()
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFilePath, nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return ret.Key
}

func setQiNiuArea(area string) *storage.Region {
	var zone *storage.Region
	switch {
	case area == "z0":
		zone = &storage.ZoneHuadong
	case area == "z1":
		zone = &storage.ZoneHuabei
	case area == "z2":
		zone = &storage.ZoneHuanan
	case area == "na0":
		zone = &storage.ZoneBeimei
	case area == "as0":
		zone = &storage.ZoneXinjiapo
	}
	return zone
}

func (qiniuUtil) DeleteFileList(fileNameList []string, qiNiuConfig map[string]interface{}) bool {
	accessKey := qiNiuConfig["qiNiuAccessKey"].(string)
	secretKey := qiNiuConfig["qiNiuSecretKey"].(string)
	bucket := qiNiuConfig["qiNiuBucket"].(string)
	successCount := 0
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{}
	// 是否使用https域名
	cfg.UseHTTPS = false
	bucketManager := storage.NewBucketManager(mac, &cfg)
	for _, fileName := range fileNameList {
		key := fileName
		err := bucketManager.Delete(bucket, key)
		if err != nil {
			fmt.Println(err)
			return false
		}
		l.Print("七牛云文件删除成功")
		successCount += 1
	}
	return true
}

func (qiniuUtil) DeleteFile(fileName string, qiNiuConfig map[string]interface{}) int {
	accessKey := qiNiuConfig["qiNiuAccessKey"].(string)
	secretKey := qiNiuConfig["qiNiuSecretKey"].(string)
	bucket := qiNiuConfig["qiNiuBucket"].(string)
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{}
	// 是否使用https域名
	cfg.UseHTTPS = false
	bucketManager := storage.NewBucketManager(mac, &cfg)
	key := fileName
	err := bucketManager.Delete(bucket, key)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	l.Print("七牛云文件删除成功")
	return -1
}

var QiniuUtil = &qiniuUtil{}
