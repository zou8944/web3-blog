package models

import (
	"blog-web3/pkg/database"
	"blog-web3/pkg/logger"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Article struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username"`
	Title     string         `json:"title" gorm:"uniqueIndex"`
	Content   string         `json:"content"`
	Tags      datatypes.JSON `json:"tags"`
	Visible   string         `json:"visible"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func ListArticleByUser(username string) []Article {
	var articles []Article
	database.DB.Model(Article{}).Where("username = ?", username).Find(&articles)
	if database.DB.Error != nil {
		logger.Errorf("List article fail. %v", database.DB.Error)
		return nil
	}
	return articles
}

func GetArticleByTitle(username string, title string) *Article {
	var article Article
	database.DB.Model(Article{}).Where("username = ? and title = ?", username, title).Find(&article)
	return &article
}

func (a *Article) Create() bool {
	database.DB.Create(a)
	return a.ID > 0
}

func (a *Article) Save() bool {
	database.DB.Save(a)
	return database.DB.Error == nil
}
