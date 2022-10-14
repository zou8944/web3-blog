package model

import (
	"github.com/pkg/errors"
	"time"
)

type User struct {
	ID            int64     `json:"id" gorm:"primary key"`
	PublicAddress string    `json:"public_address" gorm:"uniqueIndex"`
	UniqueName    string    `json:"unique_name" gorm:"uniqueIndex"`
	Nonce         string    `json:"nonce"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

func GetByPublicAddress(pa string) (*User, error) {
	var user User
	db.Where("public_address = ?", pa).Find(&user)
	return &user, errors.WithStack(db.Error)
}

func (user *User) Save() (*User, error) {
	result := db.Save(user)
	return user, result.Error
}
