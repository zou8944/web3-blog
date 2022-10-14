package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UnixTime time.Time

func (t UnixTime) MarshalJSON() ([]byte, error) {
	unixTime := fmt.Sprintf("%d", time.Time(t).UnixMilli())
	return []byte(unixTime), nil
}

type R struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Response(c *gin.Context, status int, v interface{}) {
	data := R{
		Code:    "Success",
		Message: "",
		Data:    v,
	}
	c.JSON(status, data)
}

func ResponseSuccess(c *gin.Context, v interface{}) {
	data := R{
		Code:    "Success",
		Message: "",
		Data:    v,
	}
	c.JSON(http.StatusOK, data)
}

func ResponseReject(c *gin.Context, err error) {
	data := R{
		Code:    "Rejected",
		Message: err.Error(),
		Data:    nil,
	}
	c.JSON(http.StatusBadRequest, data)
}

func ResponseError(c *gin.Context, err error) {
	data := R{
		Code:    "InternalError",
		Message: err.Error(),
		Data:    nil,
	}
	c.JSON(http.StatusInternalServerError, data)
}
