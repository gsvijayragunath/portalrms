package handlers

import (
	"net/http"
	"example.com/RMS/db"
	"example.com/RMS/errors"
	"example.com/RMS/models"
	"example.com/RMS/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobHandler struct{}

func NewJobHandler() *JobHandler {
	return &JobHandler{}
}

func (h *JobHandler) CreateJob(c *gin.Context) {

	var jobOpening models.Job
	userID, val := utils.GetUserID(c, "admin")
	if !val {
		return
	}

	jobOpening.PostedByID = userID

	err := c.ShouldBindJSON(&jobOpening)
	if err != nil {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest, err.Error(), "Invalid Input")
		c.JSON(httpStatus, response)
		return
	}

	if err = db.DB.Create(&jobOpening).Error; err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Failed to Create Job")
		c.JSON(httpStatus, response)
		return
	}

	c.JSON(http.StatusOK, utils.RenderSuccess(jobOpening))
}

func (h *JobHandler) GetAllJobByUserID(c *gin.Context) { //Userid = Admin

	var jobs []models.Job
	UserID, val := utils.GetUserID(c, "admin")
	if !val {
		return
	}

	if err := db.DB.Where("posted_by_id=?", UserID).Find(&jobs).Error; err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Failed to Fetch Jobs")
		c.JSON(httpStatus, response)
		return
	}

	c.JSON(http.StatusOK, utils.RenderSuccess(jobs))
}

func (h *JobHandler) GetJobByJobID(c *gin.Context) {

	var job models.Job
	val := utils.CheckUserType(c, "admin")
	if !val {
		return
	}

	jobID := c.Param("jobID")
	if len(jobID) != 36 {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest, "Invalid Job ID")
		c.JSON(httpStatus, response)
		return
	}

	if err := db.DB.Where("job_id=?", jobID).First(&job).Error; err != nil {
		httpStatus, response := utils.RenderError(errors.ErrNotFound, err.Error())
		c.JSON(httpStatus, response)
		return
	}

	c.JSON(http.StatusOK, utils.RenderSuccess(job))
}

func (h *JobHandler) GetAllJobs(c *gin.Context) {

	var jobs []models.Job
	if err := db.DB.Find(&jobs).Error; err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Failed to Fetch Jobs")
		c.JSON(httpStatus, response)
		return
	}

	if len(jobs) == 0 {
		c.JSON(http.StatusOK, utils.RenderSuccess("No Data Available"))
		return
	}

	c.JSON(http.StatusOK, utils.RenderSuccess(jobs))
}

func (h *JobHandler) DeleteJob(c *gin.Context) {
	userID, val := utils.GetUserID(c, "admin")
	if !val {
		return
	}

	jobIDStr := c.Param("jobID")
	if len(jobIDStr) != 36 {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest, "Invalid jobID")
		c.JSON(httpStatus, response)
		return
	}

	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		httpStatus, response := utils.RenderError(errors.ErrInvalidRequest, "Invalid UUID format")
		c.JSON(httpStatus, response)
		return
	}

	var postedByID string
	if err := db.DB.Table("jobs").
		Where("job_id = ?", jobID).
		Select("posted_by_id").
		Scan(&postedByID).Error; err != nil {
		httpStatus, response := utils.RenderError(errors.ErrNotFound, err.Error(), "Job not found")
		c.JSON(httpStatus, response)
		return
	}

	if postedByID != userID.String() {
		httpStatus, response := utils.RenderError(errors.ErrForbidden, "Access Denied. Only the job creator can delete this job.")
		c.JSON(httpStatus, response)
		return
	}

	//Psql Transaction
	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Job{}, jobID).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Application{}).
			Where("job_id = ?", jobID).
			Update("job_status", "job expired").Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		httpStatus, response := utils.RenderError(err, err.Error(), "Failed to delete job")
		c.JSON(httpStatus, response)
		return
	}

	c.JSON(http.StatusOK, utils.RenderSuccess("Job deleted successfully"))
}
