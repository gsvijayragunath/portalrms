package handlers

import (
	"fmt"
	"net/http"

	"example.com/RMS/db"
	"example.com/RMS/errors"
	"example.com/RMS/models"
	"example.com/RMS/utils"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct{}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{}
}

func (h *ProfileHandler) CreateProfile(c *gin.Context) {

	userID, val := utils.GetUserID(c, "applicant")
	if !val {
		return
	}

	var profileCount int64
	if err := db.DB.Model(&models.Profile{}).Where("user_id = ?", userID).Count(&profileCount).Error; err != nil {
		httpStatus, response := utils.RenderError(errors.ErrDatabase, "Failed to count existing profiles", "Database error")
		c.JSON(httpStatus, response)
		return
	}

	if profileCount >= 5 {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest, "Profile creation limit reached", "Maximum 5 profiles allowed")
		c.JSON(httpStatus, response)
		return
	}

	var userProfile models.Profile
	if err := c.ShouldBindJSON(&userProfile); err != nil {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest, err.Error(), "Invalid input")
		c.JSON(httpStatus, response)
		return
	}

	userProfile.UserID = userID

	if err := db.DB.Create(&userProfile).Error; err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Failed to create profile")
		c.JSON(httpStatus, response)
		return
	}

	c.JSON(http.StatusOK, utils.RenderSuccess(userProfile))
}

func (h *ProfileHandler) GetAllProfilesByUserID(c *gin.Context) {

	var profiles []models.Profile
	userID, val := utils.GetUserID(c, "applicant")
	if !val {
		return
	}

	if err := db.DB.Where("user_id = ?", userID).Find(&profiles).Error; err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Failed to fetch profiles")
		c.JSON(httpStatus, response)
		return
	}

	c.JSON(http.StatusOK, utils.RenderSuccess(profiles))
}

func (h *ProfileHandler) GetProfileByProfileID(c *gin.Context) {

	var profile models.Profile
	if !utils.CheckUserType(c, "applicant") {
		return
	}

	profileID := c.Param("profileID")
	if len(profileID) != 36 {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest, "Invalid Profile ID")
		c.JSON(httpStatus, response)
		return
	}

	if err := db.DB.Where("profile_id = ?", profileID).First(&profile).Error; err != nil {
		httpStatus, response := utils.RenderError(errors.ErrNotFound, fmt.Sprintf("record not found on this profileID:%s", profileID))
		c.JSON(httpStatus, response)
		return
	}

	c.JSON(http.StatusOK, utils.RenderSuccess(profile))
}

func (h *ProfileHandler) DeleteProfile(c *gin.Context) {
	userID, val := utils.GetUserID(c, "applicant")
	if !val {
		return
	}

	profileID := c.Param("profileID")
	if len(profileID) != 36 {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest, "Invalid Profile ID")
		c.JSON(httpStatus, response)
		return
	}

	var userProfile models.Profile
	if err := db.DB.Model(&models.Profile{}).Where("profile_id = ?", profileID).First(&userProfile).Error; err != nil {
		httpStatus, response := utils.RenderError(err, "Profile not found", "No profile exists with the given ID")
		c.JSON(httpStatus, response)
		return
	}

	if userProfile.UserID != userID {
		httpStatus, response := utils.RenderError(errors.ErrForbidden,"You do not have permission to delete this profile")
		c.JSON(httpStatus, response)
		return
	}

	if err := db.DB.Delete(&userProfile).Error; err != nil {
		httpStatus, response := utils.RenderError(errors.ErrDatabase, "Failed to delete profile")
		c.JSON(httpStatus, response)
		return
	}

	c.JSON(http.StatusOK, utils.RenderSuccess("Profile deleted successfully"))
}
