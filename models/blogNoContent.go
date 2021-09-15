package models

/**
 *
 * @author  镜湖老杨
 * @date  2020/12/25 10:25 上午
 * @version 1.0
 */
import (
	_ "gorm.io/gorm"
	_type "mogu-go-v2/models/type"
)

type BlogNoContent struct {
	Uid          string       `gorm:"primaryKey" json:"uid"`
	Title        string       `json:"title"`
	Summary      string       `json:"summary"`
	TagUid       string       `json:"tagUid"`
	ClickCount   int          `json:"clickCount"`
	CollectCount int          `json:"collectCount"`
	FileUid      string       `json:"fileUid"`
	Status       int8         `gorm:"default:1" json:"status"`
	CreatedAt    _type.MyTime `gorm:"column:create_time" json:"createTime"`
	UpdatedAt    _type.MyTime `gorm:"column:update_time" json:"updateTime"`
	AdminUid     string       `json:"adminUid"`
	IsOriginal   string       `json:"isOriginal"`
	Author       string       `json:"author"`
	ArticlesPart string       `json:"articlesPart"`
	BlogSortUid  string       `json:"blogSortUid"`
	Level        int          `json:"level"`
	IsPublish    string       `json:"isPublish"`
	Sort         int          `json:"sort"`
	OpenComment  int          `json:"openComment"`
	Type         int          `json:"type"`
	OutsideLink  string       `json:"outsideLink"`
	Oid          int          `json:"oid"`
	TagList      []Tag        `gorm:"foreignKey:Uid;references:TagUid" json:"tagList"`
	PhotoList    []string     `gorm:"-" json:"photoList"`
	BlogSort     BlogSort     `gorm:"foreignKey:BlogSortUid;references:Uid" json:"blogSort"`
	BlogSortName string       `gorm:"-" json:"blogSortName"`
	PhotoUrl     string       `gorm:"-" json:"photoUrl"`
	ParseCount   int          `gorm:"-" json:"parseCount"`
	Copyright    string       `gorm:"-" json:"copyright"`
}

func (*BlogNoContent) TableName() string {
	return "t_blog"
}
