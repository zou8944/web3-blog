package controller

import (
	"blog-web3/app/models"
	"blog-web3/pkg/jwt"
	"blog-web3/pkg/response"
	"blog-web3/pkg/web3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var inputUser models.User
	if err := c.ShouldBind(inputUser); err != nil {
		response.AbortWith400(c, err)
		return
	}
	uid, err := uuid.NewRandom()
	if err != nil {
		response.AbortWith500(c, err)
		return
	}
	inputUser.Nonce = uid.String()
	if _, err := inputUser.Save(); err != nil {
		response.AbortWith500(c, err)
		return
	}
	response.Data(c, nil)
}

func (uc *UserController) OverrideUser(c *gin.Context) {
	var inputUser models.User
	if err := c.ShouldBind(inputUser); err != nil {
		response.AbortWith400(c, err)
		return
	}
	publicAddress := c.Param("publicAddress")
	existUser, err := models.GetByPublicAddress(publicAddress)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.StatusData(c, http.StatusNotFound, nil)
		return
	}
	existUser.UniqueName = inputUser.UniqueName
	newUser, err := existUser.Save()
	if err != nil {
		response.AbortWith500(c, err)
		return
	}
	response.Data(c, newUser)
}

func (uc *UserController) GetUser(c *gin.Context) {
	publicAddress := c.Param("publicAddress")
	user, err := models.GetByPublicAddress(publicAddress)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.StatusData(c, http.StatusNotFound, nil)
		return
	}
	response.Data(c, user)
}

func (uc *UserController) LoginWithMetaMask(c *gin.Context) {
	var inputBody map[string]string
	if err := c.ShouldBind(&inputBody); err != nil {
		response.AbortWith400(c, err)
	}
	publicAddress := inputBody["publicAddress"]
	signature := inputBody["signature"]
	user, err := models.GetByPublicAddress(publicAddress)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.StatusData(c, http.StatusNotFound, nil)
		return
	}
	nonce := user.Nonce

	if sigValid := web3.VerifySignature(publicAddress, signature, nonce); !sigValid {
		response.AbortWith400(c, errors.New("signature invalid"))
		return
	}
	if _jwt := jwt.GenerateJWT(user); _jwt != "" {
		response.Data(c, gin.H{
			"token": _jwt,
			"user":  user,
		})
	} else {
		response.AbortWith500(c, err)
	}
}
