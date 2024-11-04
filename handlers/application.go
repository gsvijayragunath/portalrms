package handlers

import (
	"example.com/RMS/db"
	"example.com/RMS/errors"
	"example.com/RMS/models"
	"example.com/RMS/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type ApplicationHandler struct{}

func NewApplicationHandler() *ApplicationHandler {
	return &ApplicationHandler{}
}

type Status struct {
	Status string `json:"status" binding:"required"`
}

func (h *ApplicationHandler) FillApplication(c *gin.Context) {

	var application models.Application
	userID, val := utils.GetUserID(c, "applicant")
	if !val {
		return
	}

	jobIDstr := c.Param("jobID")
	if len(jobIDstr) != 36 {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest, "Invalid Job ID")
		c.JSON(httpStatus, response)
		return
	}

	jobID, err := uuid.Parse(jobIDstr)
	if err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Error During JobID conversion! Try Again")
		c.JSON(httpStatus, response)
		return
	}

	application.UserID = userID
	application.JobID = jobID
	application.Status = "Applied"   //Status Changeable 
	application.JobStatus = "Active" //Default -> Job Expired(Autochange when Job deleted)

	if err := c.ShouldBindJSON(&application); err != nil {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest, err.Error(), "Invalid Input")
		c.JSON(httpStatus, response)
		return
	}

	var existingApplication models.Application
	if err := db.DB.Where("user_id = ? AND job_id = ?", userID, jobID).First(&existingApplication).Error; err == nil {
		httpStatus, response := utils.RenderError(errors.ErrConflict, "Application already exists for this job")
		c.JSON(httpStatus, response)
		return
	}

	if err := db.DB.Create(&application).Error; err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Failed to Submit Application")
		c.JSON(httpStatus, response)
		return
	}
	c.JSON(http.StatusOK, utils.RenderSuccess("Applied Successfully"))
}

func (h *ApplicationHandler) GetAllApplicationsByUserID(c *gin.Context) {

	userID, val := utils.GetUserID(c, "applicant")
	if !val {
		return
	}

	var applications []models.Application
	if err := db.DB.Where("user_id = ?", userID).Find(&applications).Error; err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Failed to Fetch Applications")
		c.JSON(httpStatus, response)
		return
	}

	var profileIDs []uuid.UUID
	var jobIDs []uuid.UUID

	for _, application := range applications {
		profileIDs = append(profileIDs, application.ProfileID)
		jobIDs = append(jobIDs, application.JobID)
	}

	var profiles []models.Profile
	if err := db.DB.Where("profile_id IN ?", profileIDs).Find(&profiles).Error; err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Failed to Fetch Profiles")
		c.JSON(httpStatus, response)
		return
	}

	var jobs []models.Job
	if err := db.DB.Where("job_id IN ?", jobIDs).Find(&jobs).Error; err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Failed to Fetch Jobs")
		c.JSON(httpStatus, response)
		return
	}

	profileMap := make(map[uuid.UUID]models.Profile)
	for _, profile := range profiles {
		profileMap[profile.ProfileID] = profile
	}

	jobMap := make(map[uuid.UUID]models.Job)
	for _, job := range jobs {
		jobMap[job.JobID] = job
	}

	for i, application := range applications {
		if profile, exists := profileMap[application.ProfileID]; exists {
			applications[i].Profile = &profile
		} else {
			applications[i].Profile = nil
		}

		if job, exists := jobMap[application.JobID]; exists {
			applications[i].Job = &job
		} else {
			applications[i].Job = nil
		}
	}

	c.JSON(http.StatusOK, utils.RenderSuccess(applications))
}

func (h *ApplicationHandler) GetApplicantsByJobID(c *gin.Context) {
	val := utils.CheckUserType(c, "admin")
	if !val {
		return
	}
	jobID := c.Param("jobID")
	if len(jobID) != 36 {
		httpStatus, respone := utils.RenderError(errors.ErrInvalidRequest, "Invalid JobID")
		c.JSON(httpStatus, respone)
		return
	}

	var applications []models.Application
	if err := db.DB.Where("job_id = ?", jobID).Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Fetch Applications"})
		return
	}

	var profileIDs []uuid.UUID
	for _, application := range applications {
		profileIDs = append(profileIDs, application.ProfileID)
	}

	var profiles []models.Profile
	if err := db.DB.Where("profile_id IN ?", profileIDs).Find(&profiles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Fetch Profiles"})
		return
	}

	profileMap := make(map[uuid.UUID]models.Profile)
	for _, profile := range profiles {
		profileMap[profile.ProfileID] = profile
	}

	for i, application := range applications {
		if profile, exists := profileMap[application.ProfileID]; exists {
			applications[i].Profile = &profile
		} else {
			applications[i].Profile = nil
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": applications})
}

func (h *ApplicationHandler) ChangeStatus(c *gin.Context) {
	var status Status
	userID, val := utils.GetUserID(c, "admin")
	if !val {
		return
	}

	applicationIDstr := c.Param("applicationID")
	if len(applicationIDstr) != 36 {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest, "Invalid applicationID")
		c.JSON(httpStatus, response)
		return
	}

	applicationID, err := uuid.Parse(applicationIDstr)
	if err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Error during ApplicationID conversion")
		c.JSON(httpStatus, response)
		return
	}

	var postedByID string
	if err := db.DB.Table("applications").
		Select("jobs.posted_by_id").
		Joins("JOIN jobs ON jobs.job_id  = applications.job_id").
		Where("applications.application_id = ?", applicationID).
		Scan(&postedByID).Error; err != nil {
		httpStatus, response := utils.RenderError(errors.ErrNotFound, "Job or Application not found")
		c.JSON(httpStatus, response)
		return
	}
	// fmt.Println(postedByID)
	if postedByID != userID.String() {
		httpStatus, response := utils.RenderError(errors.ErrForbidden, "You dont have access to change the status")
		c.JSON(httpStatus, response)
		return
	}

	if err := c.ShouldBindJSON(&status); err != nil {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest, err.Error(), "Invalid Input")
		c.JSON(httpStatus, response)
		return
	}

	var application models.Application
	if err := db.DB.Where("application_id = ?", applicationID).First(&application).Error; err != nil {
		httpStatus, response := utils.RenderError(errors.ErrNotFound, "Application not found")
		c.JSON(httpStatus, response)
		return
	}

	application.Status = status.Status
	if err := db.DB.Save(&application).Error; err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Could not update status")
		c.JSON(httpStatus, response)
		return
	}

	c.JSON(http.StatusOK, utils.RenderSuccess("Status Updated Successfully"))
}

func (h *ApplicationHandler) DeleteApplication(c *gin.Context) {
	if !utils.CheckUserType(c, "applicant") {
		return
	}

	userID, val := utils.GetUserID(c, "applicant")
	if !val {
		return
	}

	applicationIDstr := c.Param("applicationID")
	if len(applicationIDstr) != 36 {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest, "Invalid applicationID")
		c.JSON(httpStatus, response)
		return
	}

	var userIDCheck string
	if err := db.DB.Model(&models.Application{}).Select("user_id").Where("application_id = ?", applicationIDstr).Scan(&userIDCheck).Error; err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Unable to fetch Required details to process with this operation")
		c.JSON(httpStatus, response)
		return
	}
	if userID.String() != userIDCheck {
		httpStatus, response := utils.RenderError(errors.ErrForbidden, "Only the owner of the application should able to delete the job", "Access Denied to delete this Application")
		c.JSON(httpStatus, response)
		return
	}

	applicationID, err := uuid.Parse(applicationIDstr)
	if err != nil {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest, "Invalid UUID format")
		c.JSON(httpStatus, response)
		return
	}

	if err := db.DB.Where("application_id = ?", applicationID).Delete(&models.Application{}).Error; err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Failed to delete application")
		c.JSON(httpStatus, response)
		return
	}

	c.JSON(http.StatusOK, utils.RenderSuccess("Application deleted successfully"))
}
