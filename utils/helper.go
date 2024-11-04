package utils

import (
	CustomError "errors"
	"example.com/RMS/errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserID(c *gin.Context, allowedUserType string) (uuid.UUID, bool) {

	userType, exists := c.Get("userType")
	if !exists {
		httpStatus, response := RenderError(CustomError.New("err"), "Error in JWT Token service")
		c.JSON(httpStatus, response)
		return uuid.Nil, false
	}

	userTypeStr, ok := userType.(string)
	if !ok {
		httpStatus, response := RenderError(CustomError.New("err"), "Error in Assertion Check :UserType")
		c.JSON(httpStatus, response)
		return uuid.Nil, false
	}

	if userTypeStr != allowedUserType {
		httpStatus, response := RenderError(errors.ErrForbidden, fmt.Sprintf("Access denied for user type: %s", userTypeStr), "Access Denied")
		c.JSON(httpStatus, response)
		return uuid.Nil, false
	}

	userID, exists := c.Get("userID")
	if !exists {
		httpStatus, response := RenderError(CustomError.New("err"), "Error in JWT Token service")
		c.JSON(httpStatus, response)
		return uuid.Nil, false
	}

	userIDStr, ok := userID.(string)
	if !ok {
		httpStatus, response := RenderError(CustomError.New("err"), "Error in Assertion Check :UserID")
		c.JSON(httpStatus, response)
		return uuid.Nil, false
	}

	parsedUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		httpStatus, response := RenderError(err, "Failed During UUID type-Conversion Try Again")
		c.JSON(httpStatus, response)
		return uuid.Nil, false
	}
	return parsedUserID, true
}

func CheckUserType(c *gin.Context, allowedUserType string) bool {
	userType, exists := c.Get("userType")
	if !exists {
		httpStatus, response := RenderError(CustomError.New("err"), "Error in JWT Token service")
		c.JSON(httpStatus, response)
		return false
	}

	userTypeStr, ok := userType.(string)
	if !ok {
		httpStatus, response := RenderError(CustomError.New("err"), "Error in Assertion Check :UserType")
		c.JSON(httpStatus, response)
		return false
	}

	if userTypeStr != allowedUserType {
		httpStatus, response := RenderError(errors.ErrForbidden, fmt.Sprintf("Access denied for user type: %s", userTypeStr), "Access Denied")
		c.JSON(httpStatus, response)
		return false
	}
	return true
}
