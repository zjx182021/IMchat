package models

import (
	"TM_chat/utils"
	"fmt"

	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name    string
	OwnerId uint
	Img     string
	Desc    string
}

func CreateCommunity(c *Community) int {
	if len(c.Name) == 0 {
		fmt.Println("名称为空")
		return -1
	}
	if c.OwnerId == 0 {
		fmt.Println("拥有者ID为空")
		return -1
	}

	if err := utils.DB.Create(c).Error; err != nil {
		fmt.Println("创建失败:", err)
		return -1
	}
	return 0
}

func LoadCommunity(ownerId uint) []*Community {
	data := make([]*Community, 10)
	if err := utils.DB.Where("owner_Id = ?", ownerId).Find(&data).Error; err != nil {
		fmt.Println("查询失败:", err)
		return nil
	}
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}
