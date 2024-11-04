package routes

import (
	"example.com/RMS/handlers"
	"example.com/RMS/middlewares"
	"github.com/gin-gonic/gin"
)

func Routes(server *gin.Engine) {

	profileHandler := handlers.NewProfileHandler()
	authHandler := handlers.NewAuthHandler()
	jobHandler := handlers.NewJobHandler()
	applicationHandler := handlers.NewApplicationHandler()

	//Authenticated routes
	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)

	//Auth
	server.POST("/signup", authHandler.Signup)
	server.POST("/signin", authHandler.Signin)

	//Profile
	authenticated.POST("/uploadresume", profileHandler.CreateProfile)
	authenticated.GET("/applicant/profiles", profileHandler.GetAllProfilesByUserID)
	authenticated.GET("/applicant/profile/:profileID", profileHandler.GetProfileByProfileID)
	authenticated.DELETE("/deleteprofile/:profileID", profileHandler.DeleteProfile)

	//Job
	authenticated.GET("/jobs", jobHandler.GetAllJobs) // Accessible for Admin & Applicants
	authenticated.POST("/admin/createjob", jobHandler.CreateJob)
	authenticated.GET("/admin/job/:jobID", jobHandler.GetJobByJobID)
	authenticated.GET("/admin/job/userid", jobHandler.GetAllJobByUserID)
	authenticated.DELETE("/admin/deletejob/:jobID", jobHandler.DeleteJob)

	//Application
	authenticated.POST("/applicant/applyjob/:jobID", applicationHandler.FillApplication)
	authenticated.GET("/applicant/all-applications", applicationHandler.GetAllApplicationsByUserID) //Applicant
	authenticated.GET("/admin/all-applicants/:jobID", applicationHandler.GetApplicantsByJobID)      // Admin
	authenticated.PUT("/admin/changestatus/:applicationID", applicationHandler.ChangeStatus)
	authenticated.DELETE("/applicant/deleteapplication/:applicationID", applicationHandler.DeleteApplication)
}
