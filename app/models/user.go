package models

import (
	"blog-web3/pkg/database"
	"blog-web3/pkg/types"
	"github.com/pkg/errors"
)

type User struct {
	ID            int64          `json:"id" gorm:"primary key"`
	PublicAddress string         `json:"public_address" gorm:"uniqueIndex"`
	UniqueName    string         `json:"unique_name" gorm:"uniqueIndex"`
	Nonce         string         `json:"nonce"`
	CreatedAt     types.UnixTime `json:"-"`
	UpdatedAt     types.UnixTime `json:"-"`
}

func GetByPublicAddress(pa string) (*User, error) {
	var user User
	database.DB.Where("public_address = ?", pa).Find(&user)
	return &user, errors.WithStack(database.DB.Error)
}

func (user *User) Save() (*User, error) {
	result := database.DB.Save(user)
	return user, result.Error
}
