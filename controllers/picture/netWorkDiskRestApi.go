package picture

import (
	"encoding/json"
	"github.com/rs/xid"
	"mogu-go-v2/common"
	"mogu-go-v2/controllers/base"
	"mogu-go-v2/models"
	"mogu-go-v2/models/vo"
	"mogu-go-v2/service"
	"strings"
)

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/13 2:36 下午
 * @version 1.0
 */

type NetWorkDiskRestApi struct {
	base.BaseController
}

var treeId int64

func (c *NetWorkDiskRestApi) GetFileList() {
	c.CheckLogin()
	var networkDisk models.NetworkDisk
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &networkDisk)
	if err != nil {
		panic(err)
	}
	networkDisk.FilePath = common.PathUtil.UrlDecode(networkDisk.FilePath)
	//token := c.Ctx.GetCookie("Admin-Token")
	header := c.Ctx.Request.Header
	token := header.Get("Authorization")
	qiNiuResultMap := c.GetSystemConfigMap(token)
	picturePriority := qiNiuResultMap["picturePriority"].(string)
	where := "status=?"
	if networkDisk.FileType != 0 {
		if networkDisk.FileType == 5 {
			where += " and extend_name not in ('" + strings.Join(common.FileUtil.GetFileExtendsByType(5), "','") + "')"
		} else {
			where += " and extend_name in ('" + strings.Join(common.FileUtil.GetFileExtendsByType(networkDisk.FileType), "','") + "')"
		}
	} else if networkDisk.FilePath != "" {
		where += " and file_path='" + networkDisk.FilePath + "'"
	}
	var list []models.NetworkDisk
	common.DB.Where(where, 1).Find(&list)
	for i, item := range list {
		if picturePriority == "1" {
			list[i].FileUrl = qiNiuResultMap["qiNiuPictureBaseUrl"].(string) + item.QiNiuUrl
		} else if picturePriority == "2" {
			list[i].FileUrl = qiNiuResultMap["minioPictureBaseUrl"].(string) + item.MinioUrl
		} else {
			list[i].FileUrl = qiNiuResultMap["localPictureBaseUrl"].(string) + item.LocalUrl
		}
	}
	c.SuccessWithData(list)
}

func (c *NetWorkDiskRestApi) CreateFile() {
	adminUid := c.CheckLogin()
	var networkDisk models.NetworkDisk
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &networkDisk)
	if err != nil {
		panic(err)
	}
	networkDisk.AdminUid = adminUid
	networkDisk.Uid = xid.New().String()
	common.DB.Create(&networkDisk)
	c.SuccessWithMessage("插入成功")
}

func (c *NetWorkDiskRestApi) Edit() {
	c.CheckLogin()
	var networkDiskVO vo.NetworkDiskVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &networkDiskVO)
	if err != nil {
		panic(err)
	}
	service.NetworkDiskService.UpdateFilepathByFilePath(networkDiskVO, true)
	c.SuccessWithMessage("更新成功")
}

func (c *NetWorkDiskRestApi) MoveFile() {
	c.CheckLogin()
	var networkDiskVO vo.NetworkDiskVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &networkDiskVO)
	if err != nil {
		panic(err)
	}
	b := service.NetworkDiskService.UpdateFilepathByFilePath(networkDiskVO, false)
	if b {
		c.SuccessWithMessage("操作成功")
	} else {
		c.ThrowError("000000", "不能选择自己！")
	}
}

func (c *NetWorkDiskRestApi) DeleteFile() {
	c.CheckLogin()
	var networkDiskVO vo.NetworkDiskVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &networkDiskVO)
	if err != nil {
		panic(err)
	}
	header := c.Ctx.Request.Header
	tokenData := header.Get("Authorization")
	qiNiuConfig := c.GetSystemConfigMap(tokenData)
	c.DeleteFileService(networkDiskVO, qiNiuConfig)
	c.SuccessWithMessage("删除成功")
}

func (c *NetWorkDiskRestApi) GetFileTree() {
	var filePathList []models.NetworkDisk
	common.DB.Where("status=? and is_dir=?", 1, 1).Find(&filePathList)
	var resultTreeNode models.TreeNode
	resultTreeNode.Label = "/"
	m := map[string]string{}
	m["filepath"] = "/"
	resultTreeNode.Attributes = m
	for i := 0; i < len(filePathList); i++ {
		filePath := filePathList[i].FilePath + filePathList[i].FileName + "/"
		var queue common.Queue
		strArr := strings.Split(filePath, "/")
		for j := 0; j < len(strArr); j++ {
			if strArr[j] != "" {
				queue.Push(strArr[j])
			}
		}
		if queue.Len() == 0 {
			continue
		}
		resultTreeNode = insertTreeNode(resultTreeNode, "/", queue)
	}
	c.SuccessWithData(resultTreeNode)
}

func insertTreeNode(treeNode models.TreeNode, filepath string, nodeNameQueue common.Queue) models.TreeNode {
	childrenTreeNodes := treeNode.Children
	currentNodeName := nodeNameQueue.Peek()
	if currentNodeName == "" {
		return treeNode
	}
	m := map[string]string{}
	filepath = filepath + currentNodeName.(string) + "/"
	m["filepath"] = filepath
	if !isExistPath(childrenTreeNodes, currentNodeName.(string)) {
		var resultTreeNode models.TreeNode
		resultTreeNode.Attributes = m
		resultTreeNode.Label = nodeNameQueue.Pop().(string)
		treeId++
		resultTreeNode.Id = treeId
		childrenTreeNodes = append(childrenTreeNodes, resultTreeNode)
	} else {
		nodeNameQueue.Pop()
	}
	if nodeNameQueue.Len() != 0 {
		for i := 0; i < len(childrenTreeNodes); i++ {
			childrenTreeNode := childrenTreeNodes[i]
			if currentNodeName == childrenTreeNode.Label {
				childrenTreeNode = insertTreeNode(childrenTreeNode, filepath, nodeNameQueue)
				childrenTreeNodes = append(childrenTreeNodes[:i], childrenTreeNodes[i+1:]...)
				childrenTreeNodes = append(childrenTreeNodes, childrenTreeNode)
				treeNode.Children = childrenTreeNodes
			}
		}
	} else {
		treeNode.Children = childrenTreeNodes
	}
	return treeNode
}

func isExistPath(childrenTreeNodes []models.TreeNode, path string) bool {
	isExistPath := false
	for i := 0; i < len(childrenTreeNodes); i++ {
		if path == childrenTreeNodes[i].Label {
			isExistPath = true
		}
	}
	return isExistPath
}

func (c *NetWorkDiskRestApi) BatchMoveFile() {
	c.CheckLogin()
	var networkDiskVO vo.NetworkDiskVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &networkDiskVO)
	if err != nil {
		panic(err)
	}
	files := networkDiskVO.Files
	newFilePath := networkDiskVO.NewFilePath
	var fileLsit []vo.NetworkDiskVO
	err1 := json.Unmarshal([]byte(files), &fileLsit)
	if err1 != nil {
		panic(err1)
	}
	var b bool
	for _, file := range fileLsit {
		file.NewFilePath = newFilePath
		file.OldFilePath = file.FilePath
		b = service.NetworkDiskService.UpdateFilepathByFilePath(file, false)
		if !b {
			break
		}
	}
	if b {
		c.SuccessWithMessage("操作成功")
	} else {
		c.ThrowError("000000", "不能选择自己！")
	}
}

func (c *NetWorkDiskRestApi) BatchDeleteFile() {
	c.CheckLogin()
	var networkDiskVOList []vo.NetworkDiskVO
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &networkDiskVOList)
	if err != nil {
		panic(err)
	}
	header := c.Ctx.Request.Header
	tokenData := header.Get("Authorization")
	qiNiuConfig := c.GetSystemConfigMap(tokenData)
	for _, networkDiskVO := range networkDiskVOList {
		c.DeleteFileService(networkDiskVO, qiNiuConfig)
	}
	c.SuccessWithMessage("批量删除成功")
}
