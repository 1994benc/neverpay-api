package user

import "github.com/jinzhu/gorm"

type UserNoPasswordModel struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `gorm:"unique" json:"email"`
	Role  string `json:"role"`
}

func (userNoPassword *UserNoPasswordModel) fromUser(u *UserModel) {
	userNoPassword.ID = u.ID
	userNoPassword.CreatedAt = u.CreatedAt
	userNoPassword.DeletedAt = u.DeletedAt
	userNoPassword.Email = u.Email
	userNoPassword.Model = u.Model
	userNoPassword.Name = u.Name
	userNoPassword.Role = u.Role
	userNoPassword.UpdatedAt = u.UpdatedAt
}
