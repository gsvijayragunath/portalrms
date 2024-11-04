package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"user_id"`
	Name            string    `json:"name" binding:"required"`
	Email           string    `gorm:"uniqueIndex;not null" json:"email" binding:"required"`
	Address         string    `json:"address"`
	UserType        string    `json:"user_type" binding:"required"`
	Password        string    `json:"password" binding:"required"`
	ProfileHeadline string    `json:"profile_headline"`
}

type Profile struct {
	ProfileID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"profile_id"`
	UserID            uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ResumeFileAddress string    `json:"resume_file_address"`
	Skills            string    `json:"skills" binding:"required"`
	Education         string    `json:"education" binding:"required"`
	Experience        string    `json:"experience"`
	Name              string    `json:"name" binding:"required"`
	Email             string    `json:"email" binding:"required"`
	Phone             int       `json:"phone_number"`
	Link              string    `json:"links"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Job struct {
	JobID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"job_id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	PostedOn    time.Time `json:"posted_on" gorm:"autoUpdateTime"`
	PostedByID  uuid.UUID `gorm:"type:uuid" json:"posted_by_id"`
	CompanyName string    `json:"company_name" binding:"required"`
}

type Application struct {
	ApplicationID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"application_id"`
	UserID        uuid.UUID `gorm:"type:uuid" json:"user_id"`
	JobID         uuid.UUID `gorm:"type:uuid;not null" json:"job_id" binding:"required"`
	JobStatus     string    `json:"job_status"`
	Job           *Job      `gorm:"foreignKey:JobID" json:"job"`
	ProfileID     uuid.UUID `gorm:"type:uuid;not null" json:"profile_id" binding:"required"`
	Profile       *Profile  `gorm:"foreignKey:ProfileID" json:"profile"`
	AppliedOn     time.Time `gorm:"autoCreateTime" json:"applied_on"`
	Status        string    `json:"status" binding:"required"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Signin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
