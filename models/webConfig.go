package models

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/18 1:04 下午
 * @version 1.0
 */

import (
	_ "gorm.io/gorm"
	"time"
)

type WebConfig struct {
	Uid                  string    `gorm:"primaryKey" json:"uid"`
	Status               int8      `gorm:"default:1" json:"status"`
	CreatedAt            time.Time `gorm:"column:create_time" json:"createTime"`
	UpdatedAt            time.Time `gorm:"column:update_time" json:"updateTime"`
	Logo                 string    `json:"logo"`
	Name                 string    `json:"name"`
	Summary              string    `json:"summary"`
	Keyword              string    `json:"keyword"`
	Author               string    `json:"author"`
	RecordNum            string    `json:"recordNum"`
	Title                string    `json:"title"`
	AliPay               string    `json:"aliPay"`
	WeixinPay            string    `json:"weixinPay"`
	OpenComment          string    `json:"openComment"`
	OpenMobileComment    string    `json:"openMobileComment"`
	OpenAdmiration       string    `json:"openAdmiration"`
	OpenMobileAdmiration string    `json:"openMobileAdmiration"`
	Github               string    `json:"github"`
	Gitee                string    `json:"gitee"`
	QqNumber             string    `json:"qqNumber"`
	QqGroup              string    `json:"qqGroup"`
	WeChat               string    `json:"weChat"`
	Email                string    `json:"email"`
	ShowList             string    `json:"showList"`
	LoginTypeList        string    `json:"loginTypeList"`
	PhotoList            []string  `gorm:"-" json:"photoList"`
	LogoPhoto            string    `gorm:"-" json:"logoPhoto"`
	AliPayPhoto          string    `gorm:"-" json:"aliPayPhoto"`
	WeixinPayPhoto       string    `gorm:"-" json:"weixinPayPhoto"`
}

func (WebConfig) TableName() string {
	return "t_web_config"
}
