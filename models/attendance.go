package models

import "../models/response"

// Attendance class
type Attendance struct {
	CommonModelFields
	Student response.User `json:"student"`
	StudentID uint `gorm:"not null" json:"student_id"`
	Schedule Schedule `json:"schedule"`
	ScheduleID uint `gorm:"unique;not null" json:"schedule_id"`
	IsAttend bool `json:"is_attend"`
	Description string `json:"description"`
}