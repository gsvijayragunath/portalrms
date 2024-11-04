package middlewares

import (
	"example.com/RMS/errors"
	"example.com/RMS/services"
	"example.com/RMS/utils"
	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	token := c.Request.Header.Get("token")

	if token == "" {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest,"Token is required and cannot be empty")
		c.AbortWithStatusJSON(httpStatus, response)
		return
	}

	userID, userType, err := services.ValidateToken(token)
	if err != nil {
		httpStatus, response := utils.RenderError(errors.ErrUnauthorized, err.Error(), "Invalid Token")
		c.AbortWithStatusJSON(httpStatus, response)
		return
	}
	c.Set("userID", userID)
	c.Set("userType", userType)
	c.Next()
}
