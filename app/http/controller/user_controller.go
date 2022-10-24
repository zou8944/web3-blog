package controller

import (
	"blog-web3/app/http/requests"
	"blog-web3/app/models"
	"blog-web3/pkg/helpers"
	"blog-web3/pkg/jwt"
	"blog-web3/pkg/response"
	"blog-web3/pkg/web3"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var request requests.CreateUserRequest
	if ok := requests.BindAndValidate(c, &request); !ok {
		return
	}

	nonce := helpers.GenerateNonce()
	if nonce == "" {
		response.AbortWith500(c, errors.New("Create User fail"))
		return
	}

	user := models.User{
		PublicAddress: request.PublicAddress,
		UniqueName:    request.UniqueName,
		Nonce:         nonce,
	}

	if _, err := user.Save(); err != nil {
		response.AbortWith500(c, err)
		return
	}
	response.Data(c, user)
}

func (uc *UserController) OverrideUser(c *gin.Context) {
	var request requests.UpdateUserRequest
	if ok := requests.BindAndValidate(c, &request); !ok {
		return
	}
	publicAddress := c.Param("publicAddress")

	existUser, err := models.GetUserByPublicAddress(publicAddress)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.StatusData(c, http.StatusNotFound, nil)
		return
	}
	existUser.UniqueName = request.UniqueName
	newUser, err := existUser.Save()
	if err != nil {
		response.AbortWith500(c, err)
		return
	}
	response.Data(c, newUser)
}

func (uc *UserController) GetUser(c *gin.Context) {
	publicAddress := c.Param("publicAddress")
	user, err := models.GetUserByPublicAddress(publicAddress)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.StatusData(c, http.StatusNotFound, nil)
		return
	}
	response.Data(c, user)
}

func (uc *UserController) LoginWithMetaMask(c *gin.Context) {
	var body requests.LoginMetaMaskRequest
	if ok := requests.BindAndValidate(c, &body); !ok {
		return
	}

	user, err := models.GetUserByPublicAddress(body.PublicAddress)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.StatusData(c, http.StatusNotFound, nil)
		return
	}

	if sigValid := web3.VerifySignature(body.PublicAddress, body.Signature, user.Nonce); !sigValid {
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
