package models

import "golang.org/x/crypto/bcrypt"

// User class
type User struct {
	CommonModelFields
	Role Role `json:"role"`
	RoleID uint `gorm:"not null" json:"role_id"`
	UserID string `gorm:"not null;type:varchar(15)" json:"user_id"`
	Name string `gorm:"not null" json:"name"`
	Username string `gorm:"type:varchar(100);unique_index" json:"username"`
	Email string `gorm:"not null;type:varchar(100);unique_index" json:"email"`
	Password []byte `gorm:"not null" json:"password"`
	Schedules []Schedule `gorm:"foreignkey:LecturerID"`
}

// GeneratePassword using bcrypt
func GeneratePassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// ComparePassword hashed with bcrypt
func ComparePassword(hashedPassword, givenPassword []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, givenPassword)
	return err == nil
}

// SetPassword -> hash password
func (u *User) SetPassword(password string) error {
	hashed, err := GeneratePassword([]byte(password))
	if err != nil {
		return err
	}
	u.Password = hashed
	return nil
}

// CheckPassword -> checking password
func (u *User) CheckPassword(password string) bool {
	return ComparePassword(u.Password, []byte(password))
}