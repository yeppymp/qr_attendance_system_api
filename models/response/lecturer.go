package response

type User struct {
	ID uint `gorm:"primary_key" json:"id"`
	RoleID uint `gorm:"not null" json:"role_id"`
	UserID string `gorm:"not null;type:varchar(15)" json:"user_id"`
	Name string `gorm:"not null" json:"name"`
	Email string `gorm:"not null;type:varchar(100);unique_index" json:"email"`
}