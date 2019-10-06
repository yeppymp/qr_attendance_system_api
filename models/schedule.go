package models

import "main.go/models/response"

// Schedule class
type Schedule struct {
	CommonModelFields
	Lecturer response.User `json:"lecturer"`
	LecturerID uint `gorm:"not null" json:"lecturer_id"`
	Day string `gorm:"type:varchar(10);not null" json:"day"`
	TimeStart string `gorm:"type:time;not null" json:"time_start"`
	TimeEnd string `gorm:"type:time;not null" json:"time_end"`
	CourseName string `gorm:"type:varchar(100); not null" json:"course_name"`
	ClassName string `gorm:"type:varchar(100); not null" json:"class_name"`
}