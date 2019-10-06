package models

// Role class
type Role struct {
	CommonModelFields
	RoleName string `gorm:"not null" json:"role_name"`
	Users []User `gorm:"foreignkey:RoleID" json:"users"`
}