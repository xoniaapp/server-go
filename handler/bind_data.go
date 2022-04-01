package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xoniaapp/server/model"
	"github.com/xoniaapp/server/model/apperrors"
)

type Request interface {
	validate() error
}

func bindData(c *gin.Context, req Request) bool {
	if err := c.ShouldBind(req); err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return false
	}

	if err := req.validate(); err != nil {
		errors := strings.Split(err.Error(), ";")
		fErrors := make([]model.FieldError, 0)

		for _, e := range errors {
			split := strings.Split(e, ":")
			er := model.FieldError{
				Field:   strings.TrimSpace(split[0]),
				Message: strings.TrimSpace(split[1]),
			}
			fErrors = append(fErrors, er)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": fErrors,
		})
		return false
	}
	return true
}
