package handlers

type ErrorSchema struct {
	Error string `form:"error" json:"error" binding:"required"`
}

type LoginSchema struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"Password" json:"Password" binding:"required"`
}

type TokenSchema struct {
	Token string `form:"token" json:"token" binding:"required"`
}
