package models

/**
 *
 * @author  镜湖老杨
 * @date  2021/1/16 1:18 下午
 * @version 1.0
 */

type TreeNode struct {
	Id         int64             `json:"id"`
	Label      string            `json:"label"`
	Depth      int               `json:"depth"`
	State      string            `gorm:"default:'closed'" json:"state"`
	Attributes map[string]string `json:"attributes"`
	Children   []TreeNode        `json:"children"`
}
