package vo

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/28 10:26 上午
 * @version 1.0
 */

type SystemConfigVO struct {
	Keyword                   string      `json:"keyword"`
	CurrentPage               int         `json:"currentPage"`
	PageSize                  int         `json:"pageSize"`
	Uid                       string      `json:"uid"`
	Status                    int         `json:"status"`
	QiNiuAccessKey            string      `json:"qiNiuAccessKey"`
	QiNiuSecretKey            string      `json:"qiNiuSecretKey"`
	Email                     string      `json:"email"`
	EmailUserName             string      `json:"emailUserName"`
	EmailPassword             string      `json:"emailPassword"`
	SmtpAddress               string      `json:"smtpAddress"`
	SmtpPort                  string      `json:"smtpPort"`
	QiNiuBucket               string      `json:"qiNiuBucket"`
	QiNiuArea                 string      `json:"qiNiuArea"`
	UploadQiNiu               string      `json:"uploadQiNiu"`
	UploadLocal               string      `json:"uploadLocal"`
	PicturePriority           string      `json:"picturePriority"`
	QiNiuPictureBaseUrl       string      `json:"qiNiuPictureBaseUrl"`
	LocalPictureBaseUrl       string      `json:"localPictureBaseUrl"`
	StartEmailNotification    string      `json:"startEmailNotification"`
	EditorModel               interface{} `json:"editorModel"`
	ThemeColor                string      `json:"themeColor"`
	MinioEndPoint             string      `json:"minioEndPoint"`
	MinioAccessKey            string      `json:"minioAccessKey"`
	MinioSecretKey            string      `json:"minioSecretKey"`
	MinioBucket               string      `json:"minioBucket"`
	UploadMinio               string      `json:"uploadMinio"`
	MinioPictureBaseUrl       string      `json:"minioPictureBaseUrl"`
	OpenDashboardNotification interface{} `json:"openDashboardNotification"`
	DashboardNotification     string      `json:"dashboardNotification"`
	ContentPicturePriority    interface{} `json:"contentPicturePriority"`
}
