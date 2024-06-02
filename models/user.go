package models

import (
	"gorm.io/gorm"
)

// GetUserByName 根据username获取User对象/**
func (m *User) GetUserByName(username string) (row User, err error) {
	var o *gorm.DB
	if o, err = NewOrm(); err != nil {
		return
	}
	sqlStr := `select * from user where username=? group by t.id`
	o.Raw(sqlStr, username).First(&row)
	return
}
