package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/15 3:30 下午
 * @version 1.0
 */
type NetworkDiskVO struct {
	Uid           string `json:"uid"`
	Status        int    `json:"status"`
	AdminUid      string `json:"adminUid"`
	FileUrl       string `json:"fileUrl"`
	FilePath      string `json:"filePath"`
	FileOldName   string `json:"fileOldName"`
	TimestampName string `json:"timestampName"`
	ExtendName    string `json:"extendName"`
	FileName      string `json:"fileName"`
	FileSize      int64  `json:"fileSize"`
	IsDir         int    `json:"isDir"`
	OldFilePath   string `json:"oldFilePath"`
	NewFilePath   string `json:"newFilePath"`
	Files         string `json:"files"`
	FileType      int    `json:"fileType"`
}
