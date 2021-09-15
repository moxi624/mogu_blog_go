package routers

import (
	"github.com/beego/beego/v2/server/web"
	"mogu-go-v2/controllers/picture"
)

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/31 2:49 下午
 * @version 1.0
 */

func init() {
	// 请求前缀
	prefix := "/mogu-picture"

	file := web.NewNamespace(prefix + "/file",
		web.NSRouter("/getPicture", &picture.FileRestApi{}, "get:GetPicture"),
		web.NSRouter("/cropperPicture", &picture.FileRestApi{}, "post:CropperPicture"),
		web.NSRouter("/pictures", &picture.FileRestApi{}, "post:UploadPics"),
	)
	networkDisk := web.NewNamespace(prefix + "/networkDisk",
		web.NSRouter("/getFileList", &picture.NetWorkDiskRestApi{}, "post:GetFileList"),
		web.NSRouter("/createFile", &picture.NetWorkDiskRestApi{}, "post:CreateFile"),
		web.NSRouter("/edit", &picture.NetWorkDiskRestApi{}, "post:Edit"),
		web.NSRouter("/deleteFile", &picture.NetWorkDiskRestApi{}, "post:DeleteFile"),
		web.NSRouter("/getFileTree", &picture.NetWorkDiskRestApi{}, "post:GetFileTree"),
		web.NSRouter("/moveFile", &picture.NetWorkDiskRestApi{}, "post:MoveFile"),
		web.NSRouter("/batchMoveFile", &picture.NetWorkDiskRestApi{}, "post:BatchMoveFile"),
		web.NSRouter("/batchDeleteFile", &picture.NetWorkDiskRestApi{}, "post:BatchDeleteFile"),
	)
	storage := web.NewNamespace(prefix + "/storage",
		web.NSRouter("/getStorage", &picture.StorageRestApi{}, "get:GetStorage"),
		web.NSRouter("/uploadFile", &picture.StorageRestApi{}, "post:UploadFile"),
	)
	webfile := web.NewNamespace(prefix + "mogu-picture/file",
		web.NSRouter("/cropperPicture", &picture.FileRestApi{}, "post:CropperPicture"),
	)
	web.AddNamespace(file, networkDisk, storage, webfile)
}
