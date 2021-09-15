package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/19 4:23 下午
 * @version 1.0
 */
type PictureVO struct {
	Uid            string `json:"uid"`
	Status         int    `json:"status"`
	Keyword        string `json:"keyword"`
	CurrentPage    int    `json:"currentPage"`
	PageSize       int    `json:"pageSize"`
	FileUid        string `json:"fileUid"`
	FileUids       string `json:"fileUids"`
	PicName        string `json:"picName"`
	PictureSortUid string `json:"pictureSortUid"`
}
