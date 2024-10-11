package models

import (
	"TM_chat/utils"
	"fmt"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	OwnerId  uint
	TargetId uint
	Type     int
	Desc     string
}

func (table *Contact) ContactTableName() string {
	return "Contact_basics"
}

func SearchFriend(userId uint) []*UserBasic {
	contacts := make([]*Contact, 0)
	objId := make([]uint, 0)
	utils.DB.Where("owner_id = ? and type = ?", userId, 1).Find(&contacts)
	for _, data := range contacts {
		objId = append(objId, data.TargetId)
	}
	users := make([]*UserBasic, 0)
	utils.DB.Where("id IN (?)", objId).Find(&users)
	return users
}

func AddFriend(userId uint, targetName string) int {
	user := &UserBasic{}
	user.Name = targetName
	fmt.Println("用户:", targetName)
	_ = FindUserByname(user)
	if user.ID != 0 {
		if user.ID == userId {
			return -1
		}
		fmt.Println("userId:   ", userId, "targetId:    ", user.ID)
		contact := &Contact{}
		contact1 := &Contact{}
		contact0 := &Contact{}
		utils.DB.Where("owner_id = ? and target_id = ? and type = 1", userId, user.ID).Find(&contact0)
		if contact0.ID != 0 {
			return -1
		}
		contact.OwnerId = userId
		contact.TargetId = user.ID
		contact.Type = 1
		contact1.OwnerId = user.ID
		contact1.TargetId = userId
		contact1.Type = 1
		tx := utils.DB.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()
		if err := utils.DB.Create(contact).Error; err != nil {
			tx.Rollback()
			return -1
		}
		if err := utils.DB.Create(contact1).Error; err != nil {
			tx.Rollback()
			return -1
		}
		tx.Commit()
		return 0

	}
	fmt.Println("用户没找到")
	return -1
}
