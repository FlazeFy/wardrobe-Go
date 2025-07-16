package others

type LoginRequest struct {
	Email    string `json:"email" binding:"required,min=6,max=36,email" example:"flazen.test@gmail.com"`
	Password string `json:"password" binding:"required,min=6,max=36" example:"nopassword"`
}
