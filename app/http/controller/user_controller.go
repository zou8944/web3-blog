package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/project5e/web3-blog/app/http/requests"
	"github.com/project5e/web3-blog/app/models"
	"github.com/project5e/web3-blog/pkg/helpers"
	"github.com/project5e/web3-blog/pkg/jwt"
	"github.com/project5e/web3-blog/pkg/response"
	"github.com/project5e/web3-blog/pkg/web3"
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
		response.AbortWith500(c)
		return
	}

	user := models.User{
		PublicAddress: request.PublicAddress,
		UniqueName:    request.UniqueName,
		Nonce:         nonce,
	}

	if _, err := user.Save(); err != nil {
		response.AbortWith500(c)
		return
	}
	response.Created(c, user)
}

func (uc *UserController) OverrideUser(c *gin.Context) {
	var request requests.UpdateUserRequest
	if ok := requests.BindAndValidate(c, &request); !ok {
		return
	}
	publicAddress := c.Param("publicAddress")

	existUser := models.GetUserByPublicAddress(publicAddress)
	if existUser.ID <= 0 {
		response.NotFound(c)
		return
	}
	existUser.UniqueName = request.UniqueName
	newUser, err := existUser.Save()
	if err != nil {
		response.AbortWith500(c)
		return
	}
	response.SuccessWithData(c, newUser)
}

func (uc *UserController) GetUser(c *gin.Context) {
	publicAddress := c.Param("publicAddress")
	user := models.GetUserByPublicAddress(publicAddress)
	if user.ID <= 0 {
		response.NotFound(c)
		return
	}
	response.SuccessWithData(c, user)
}

func (uc *UserController) LoginWithMetaMask(c *gin.Context) {
	var body requests.LoginMetaMaskRequest
	if ok := requests.BindAndValidate(c, &body); !ok {
		return
	}

	user := models.GetUserByPublicAddress(body.PublicAddress)
	if user.ID <= 0 {
		response.NotFound(c)
		return
	}

	// todo move to requests package
	if sigValid := web3.VerifySignature(body.PublicAddress, body.Signature, user.Nonce); !sigValid {
		response.AbortWith400(c, errors.New("signature invalid"))
		return
	}
	user.Nonce = helpers.GenerateNonce()
	newUser, err := user.Save()
	if err != nil {
		response.AbortWith500(c)
		return
	}
	if _jwt := jwt.GenerateJWT(user); _jwt != "" {
		response.SuccessWithData(c, gin.H{
			"token": _jwt,
			"user":  newUser,
		})
	} else {
		response.AbortWith500(c)
	}
}
