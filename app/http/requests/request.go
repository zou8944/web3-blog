package requests

import (
	"blog-web3/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()
}

func BindAndValidate(c *gin.Context, value interface{}) bool {
	if err := c.ShouldBindJSON(value); err != nil {
		response.AbortWith400(c, err)
		return false
	}
	if err := validate.Struct(value); err != nil {
		errs := retrieveValidationErrors(err)
		response.StatusData(c, http.StatusBadRequest, errs)
		return false
	}
	return true
}

func retrieveValidationErrors(err error) map[string]string {
	vErrs := err.(validator.ValidationErrors)
	var errs map[string]string
	for _, vErr := range vErrs {
		errs[vErr.Field()] = vErr.Error()
	}
	return errs
}