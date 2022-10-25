package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type R struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

func Created(c *gin.Context, v interface{}) {
	c.JSON(http.StatusCreated, R{"Created", "", v, nil})
}

func Success(c *gin.Context) {
	c.JSON(http.StatusOK, R{"Success", "", nil, nil})
}

func SuccessWithData(c *gin.Context, v interface{}) {
	c.JSON(http.StatusOK, R{"Success", "", v, nil})
}

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, R{"NotFound", "Resource not found", nil, nil})
}

func AbortWith400(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, R{"BadRequest", err.Error(), nil, nil})
}

func AbortWithValidateFail(c *gin.Context, errs map[string]string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, R{"BadRequest", "Validate fail, please refer to field [errors]", nil, errs})
}

func AbortWith500(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, R{"BadRequest", "Server Internal Error", nil, nil})
}
