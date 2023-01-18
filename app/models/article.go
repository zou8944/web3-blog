package models

import (
	"github.com/goccy/go-json"
	"github.com/project5e/web3-blog/pkg/arweave"
	"github.com/project5e/web3-blog/pkg/database"
	"github.com/project5e/web3-blog/pkg/ipfs"
	"github.com/project5e/web3-blog/pkg/logger"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Article struct {
	ArWeaveTxID string         `json:"ar_weave_tx_id" gorm:"primaryKey"`
	IpfsID      string         `json:"ipfs_id"`
	Slug        string         `json:"slug" gorm:"uniqueIndex"`
	Title       string         `json:"title" gorm:"uniqueIndex"`
	Content     string         `json:"content"`
	Tags        datatypes.JSON `json:"tags"`
	Visible     string         `json:"visible"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}

func ListArticle() []Article {
	var articles []Article
	database.DB.Model(Article{}).Find(&articles)
	if database.DB.Error != nil {
		logger.Errorf("List article fail. %v", database.DB.Error)
		return nil
	}
	return articles
}

func GetArticleBySlug(slug string) *Article {
	var article Article
	if err := database.DB.Model(Article{}).Where("slug = ?", slug).First(&article).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Errorf("Get article fail. %+v", err)
		}
		return nil
	}
	return &article
}

func GetArticleByTitle(title string) *Article {
	var article Article
	database.DB.Model(Article{}).Where("title = ?", title).Find(&article)
	return &article
}

func GetArticleById(id string) *Article {
	var article Article
	err := database.DB.Model(Article{}).Where("ar_weave_tx_id = ?", id).First(&article).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Errorf("Get article fail. %+v", err)
		}
		return nil
	}
	return &article
}

func (a *Article) Create() bool {
	articleData, err := a.Web3Marshal()
	if err != nil {
		logger.ErrorIf(err)
		return false
	}
	cid, err := ipfs.UploadData(articleData)
	if err != nil {
		logger.ErrorIf(err)
		return false
	}
	a.IpfsID = cid
	txId, err := arweave.UploadData(articleData)
	if err != nil {
		logger.ErrorIf(err)
		return false
	}
	a.ArWeaveTxID = txId
	database.DB.Create(a)
	logger.ErrorIf(database.DB.Error)
	return database.DB.Error == nil
}

func (a *Article) Update() bool {
	articleData, err := a.Web3Marshal()
	if err != nil {
		logger.ErrorIf(err)
		return false
	}
	cid, err := ipfs.UploadData(articleData)
	if err != nil {
		logger.ErrorIf(err)
		return false
	}
	a.IpfsID = cid
	txId, err := arweave.UploadData(articleData)
	if err != nil {
		logger.ErrorIf(err)
		return false
	}
	a.ArWeaveTxID = txId
	database.DB.Save(a)
	logger.ErrorIf(database.DB.Error)
	return database.DB.Error == nil
}

func (a *Article) UpdateBySlug() bool {
	articleData, err := a.Web3Marshal()
	if err != nil {
		logger.ErrorIf(err)
		return false
	}
	cid, err := ipfs.UploadData(articleData)
	if err != nil {
		logger.ErrorIf(err)
		return false
	}
	a.IpfsID = cid
	txId, err := arweave.UploadData(articleData)
	if err != nil {
		logger.ErrorIf(err)
		return false
	}
	a.ArWeaveTxID = txId
	database.DB.Model(Article{}).Where("slug = ?", a.Slug).Updates(a)
	logger.ErrorIf(database.DB.Error)
	return database.DB.Error == nil
}

func (a *Article) Delete() bool {
	return a.Delete()
}

func (a *Article) Web3Marshal() ([]byte, error) {
	type _article struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	article := _article{
		Title:   a.Title,
		Content: a.Content,
	}
	return json.Marshal(article)
}
