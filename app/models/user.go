package models

import (
	"github.com/project5e/web3-blog/pkg/database"
	"time"
)

type User struct {
	ID            int64     `json:"id" gorm:"primaryKey"`
	PublicAddress string    `json:"public_address" gorm:"uniqueIndex"`
	UniqueName    string    `json:"unique_name" gorm:"uniqueIndex"`
	Nonce         string    `json:"nonce"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

func GetUserByPublicAddress(pa string) *User {
	var user User
	database.DB.Where("public_address = ?", pa).First(&user)
	return &user
}

func (user *User) Save() (*User, error) {
	result := database.DB.Save(user)
	return user, result.Error
}
