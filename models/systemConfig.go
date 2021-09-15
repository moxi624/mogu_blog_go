package models

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/21 4:36 下午
 * @version 1.0
 */
import (
	_ "gorm.io/gorm"
	"time"
)

type SystemConfig struct {
	Uid                       string    `gorm:"primaryKey" json:"uid"`
	QiNiuAccessKey            string    `json:"qiNiuAccessKey"`
	QiNiuSecretKey            string    `json:"qiNiuSecretKey"`
	Email                     string    `json:"email"`
	EmailUserName             string    `json:"emailUserName"`
	EmailPassword             string    `json:"emailPassword"`
	SmtpAddress               string    `json:"smtpAddress"`
	SmtpPort                  string    `json:"smtpPort"`
	Status                    int8      `gorm:"default:1" json:"status"`
	CreatedAt                 time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt                 time.Time `gorm:"column:update_time" json:"updateTime"`
	QiNiuBucket               string    `json:"qiNiuBucket"`
	QiNiuArea                 string    `json:"qiNiuArea"`
	UploadQiNiu               string    `json:"uploadQiNiu"`
	UploadLocal               string    `json:"uploadLocal"`
	PicturePriority           string    `json:"picturePriority"`
	QiNiuPictureBaseUrl       string    `json:"qiNiuPictureBaseUrl"`
	LocalPictureBaseUrl       string    `json:"localPictureBaseUrl"`
	StartEmailNotification    string    `json:"startEmailNotification"`
	EditorModel               string    `json:"editorModel"`
	ThemeColor                string    `json:"themeColor"`
	MinioEndPoint             string    `json:"minioEndPoint"`
	MinioAccessKey            string    `json:"minioAccessKey"`
	MinioSecretKey            string    `json:"minioSecretKey"`
	MinioBucket               string    `json:"minioBucket"`
	UploadMinio               string    `json:"uploadMinio"`
	MinioPictureBaseUrl       string    `json:"minioPictureBaseUrl"`
	OpenDashboardNotification string    `json:"openDashboardNotification"`
	DashboardNotification     string    `gorm:"type:text" json:"dashboardNotification"`
	ContentPicturePriority    string    `json:"contentPicturePriority"`
	OpenEmailActivate         string    `json:"openEmailActivate"`
	SearchModel               string    `json:"searchModel"`
}

func (SystemConfig) TableName() string {
	return "t_system_config"
}
