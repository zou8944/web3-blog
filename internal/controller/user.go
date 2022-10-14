package controller

import (
	"blog-web3/internal/infra/token"
	"blog-web3/internal/infra/web3"
	"blog-web3/internal/model"
	"blog-web3/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
)

func CreateUser(c *gin.Context) {
	var inputUser model.User
	if err := c.ShouldBind(inputUser); err != nil {
		utils.ResponseReject(c, err)
		return
	}
	uid, err := uuid.NewRandom()
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	inputUser.Nonce = uid.String()
	if _, err := inputUser.Save(); err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, nil)
}

func OverrideUser(c *gin.Context) {
	var inputUser model.User
	if err := c.ShouldBind(inputUser); err != nil {
		utils.ResponseReject(c, err)
		return
	}
	publicAddress := c.Param("publicAddress")
	existUser, err := model.GetByPublicAddress(publicAddress)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.Response(c, http.StatusNotFound, nil)
		return
	}
	existUser.UniqueName = inputUser.UniqueName
	newUser, err := existUser.Save()
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseSuccess(c, newUser)
}

func GetUser(c *gin.Context) {
	publicAddress := c.Param("publicAddress")
	user, err := model.GetByPublicAddress(publicAddress)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.Response(c, http.StatusNotFound, nil)
		return
	}
	utils.ResponseSuccess(c, user)
}

func LoginWithMetaMask(c *gin.Context) {
	var inputBody map[string]string
	if err := c.ShouldBind(&inputBody); err != nil {
		utils.ResponseReject(c, err)
	}
	publicAddress := inputBody["publicAddress"]
	signature := inputBody["signature"]
	user, err := model.GetByPublicAddress(publicAddress)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		utils.Response(c, http.StatusNotFound, nil)
		return
	}
	nonce := user.Nonce

	if sigValid := web3.VerifySignature(publicAddress, signature, nonce); !sigValid {
		utils.ResponseReject(c, errors.New("signature invalid"))
		return
	}
	jwt, err := token.GenerateJWT(user)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	type response struct {
		model.User
		Token string `json:"token"`
	}
	res := response{
		User:  *user,
		Token: jwt,
	}
	utils.ResponseSuccess(c, res)
}
