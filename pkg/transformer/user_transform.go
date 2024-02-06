package transformer

import "app-service-com/pkg/models"

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"fullname"`
}

func TransformUser(user *User, userModel models.User) {
	user.Email = userModel.Email
	user.FullName = userModel.Fullname
	user.Username = userModel.Username
}
