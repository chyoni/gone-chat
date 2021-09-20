package entity

import (
	"github.com/chiwon99881/gone-chat/utils"
	"gorm.io/gorm"
)

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	pwAsBytes := utils.ToBytes(u.Password)
	pwAsHash := utils.ToHexStringHash(pwAsBytes)
	u.Password = pwAsHash
	return
}
