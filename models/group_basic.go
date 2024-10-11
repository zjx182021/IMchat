package models

import (
	"gorm.io/gorm"
)

type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerId uint
	Icon    string
	Desc    string
	Type    int
}

func (table *GroupBasic) GroupBasicTableName() string {
	return "GroupBasic_basics"
}
