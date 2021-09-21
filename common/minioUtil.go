package common

import (
	"github.com/astaxie/beego/logs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"mogu-go-v2/models"
)

/**
 * minio上传工具类
 * @author 陌溪
 * @date  2021年9月21日08:23:03
 * @version 1.0
 */

type minioUtil struct{}

var minioClient *minio.Client

func initMinio(sysConfig models.SystemConfig)(error, *minio.Client)  {
	endpoint := sysConfig.MinioEndPoint
	accessKey := sysConfig.MinioAccessKey
	secretKey := sysConfig.MinioSecretKey
	useSSL := false
	// 初使化minio client对象。
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logs.Error("连接minio失败：error: %s", err)
		return err, minioClient
	}
	minioBucket := sysConfig.MinioBucket
	err = minioClient.MakeBucket(ctx, minioBucket, minio.MakeBucketOptions{})
	if err != nil {
		// 检查是否创建了该bucket
		exists, errBucketExists := minioClient.BucketExists(ctx, minioBucket)
		if errBucketExists == nil && exists {
			logs.Info("已经创建过Bucket: %s\n", minioBucket)
		} else {
			logs.Error(err)
		}
	} else {
		logs.Info("成功创建Bucket %s\n", minioBucket)
	}

	return nil, minioClient
}

// 上传文件到minio
func (minioUtil) UploadMinio(fileName string, localFilePath string, sysConfig models.SystemConfig) string {
	err, minioClient = initMinio(sysConfig)
	if err != nil {
		logs.Error(err)
		panic(err)
	}
	info, err := minioClient.FPutObject(ctx, sysConfig.MinioBucket, fileName, localFilePath, minio.PutObjectOptions{})
	if err != nil {
		logs.Error(err)
		panic(err)
	}
	logs.Info("文件上传成功 %s of size %d\n", fileName, info.Size)
	return info.Key
}

func (minioUtil) DeleteFileList(fileNameList []string, qiNiuConfig map[string]interface{}) bool {
	systemConfig := models.SystemConfig{
		 MinioAccessKey: qiNiuConfig["minioAccessKey"].(string),
		 MinioSecretKey: qiNiuConfig["minioNiuSecretKey"].(string),
		 MinioEndPoint: qiNiuConfig["minioEndPoint"].(string),
		 MinioBucket: qiNiuConfig["minioBucket"].(string),
	}
	err, minioClient := initMinio(systemConfig)
	if err != nil {
		logs.Error(err)
		panic(err)
	}
	for _, fileName := range fileNameList {
		err := minioClient.RemoveObject(ctx, systemConfig.MinioBucket, fileName, minio.RemoveObjectOptions{})
		if err != nil {
			logs.Error("%s 文件删除失败，error: ", fileName, err)
		} else {
			logs.Info("%s 七牛云文件删除成功", fileName)
		}
	}
	return true
}

func (minioUtil) DeleteFile(fileName string, qiNiuConfig map[string]interface{}) bool {
	systemConfig := models.SystemConfig{
		MinioAccessKey: qiNiuConfig["minioAccessKey"].(string),
		MinioSecretKey: qiNiuConfig["minioNiuSecretKey"].(string),
		MinioEndPoint: qiNiuConfig["minioEndPoint"].(string),
		MinioBucket: qiNiuConfig["minioBucket"].(string),
	}
	err, minioClient := initMinio(systemConfig)
	if err != nil {
		logs.Error(err)
		panic(err)
	}
	err = minioClient.RemoveObject(ctx, systemConfig.MinioBucket, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		logs.Error("%s 文件删除失败，error: ", fileName, err)
	} else {
		logs.Info("%s 七牛云文件删除成功", fileName)
	}
	return true
}

var MinioUtil = &minioUtil{}
