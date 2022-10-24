package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type R struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func StatusData(c *gin.Context, status int, v interface{}) {
	data := R{
		Code:    "Data",
		Message: "",
		Data:    v,
	}
	c.JSON(status, data)
}

func Data(c *gin.Context, v interface{}) {
	data := R{
		Code:    "Data",
		Message: "",
		Data:    v,
	}
	c.JSON(http.StatusOK, data)
}

func AbortWith400(c *gin.Context, err error) {
	data := R{
		Code:    "Rejected",
		Message: err.Error(),
		Data:    nil,
	}
	c.JSON(http.StatusBadRequest, data)
}

func AbortWith500(c *gin.Context, err error) {
	data := R{
		Code:    "InternalError",
		Message: err.Error(),
		Data:    nil,
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, data)
}
