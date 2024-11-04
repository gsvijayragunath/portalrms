package handlers

import (
	"example.com/RMS/db"
	"example.com/RMS/errors"
	"example.com/RMS/models"
	"example.com/RMS/services"
	"example.com/RMS/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

type AuthResponse struct {
	Token    string `json:"token"`
	UserType string `json:"user_type"`
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

func (h *AuthHandler) Signup(c *gin.Context) {

	var user models.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest, err.Error(), "Invalid Input")
		c.JSON(httpStatus, response)
		return
	}

	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Failed to Create Hashed Password")
		c.JSON(httpStatus, response)
		return
	}

	user.Password = string(hashedpassword)
	if err := db.DB.Create(&user).Error; err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Failed To Create User")
		c.JSON(httpStatus, response)
		return
	}

	c.JSON(http.StatusCreated, utils.RenderSuccess("User Created Successfully"))
}

func (h *AuthHandler) Signin(c *gin.Context) {

	var input models.Signin
	var user models.User

	err := c.ShouldBindJSON(&input)
	if err != nil {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest, err.Error(), "Invalid Input")
		c.JSON(httpStatus, response)
		return
	}

	if err := db.DB.Where("email=?", input.Email).First(&user).Error; err != nil {
		httpStatus, response := utils.RenderError(errors.ErrNotFound, err.Error(), "User Not Found")
		c.JSON(httpStatus, response)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest, err.Error(), "Password Incorrect!")
		c.JSON(httpStatus, response)
		return
	}

	token, err := services.GenerateToken(user.Email, user.UserID, user.UserType)
	if err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Failed to Generate Token. Try Again")
		c.JSON(httpStatus, response)
		return
	}

	c.JSON(http.StatusOK, utils.RenderSuccess(AuthResponse{
		Token:    token,
		UserType: user.UserType,
		UserID:   user.UserID.String(),
		UserName: user.Name,
	}))

}
