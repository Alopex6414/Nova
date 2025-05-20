package app

type UserID struct {
	UserId string `json:"userId" yaml:"userId" binding:"required"`
}

type User struct {
	UserId      string `json:"userId" yaml:"userId" binding:"required"`
	Username    string `json:"username" yaml:"username" binding:"required"`
	Password    string `json:"password" yaml:"password" binding:"required"`
	PhoneNumber string `json:"phone_number" yaml:"phone_number" binding:"required"`
	Email       string `json:"email" yaml:"email" binding:"required"`
	Address     string `json:"address" yaml:"address" binding:"required"`
	Company     string `json:"company" yaml:"company" binding:"required"`
}

type ProblemDetails struct {
	Type   string `json:"type" yaml:"type"`
	Title  string `json:"title" yaml:"title"`
	Status int    `json:"status" yaml:"status"`
	Cause  string `json:"cause" yaml:"cause"`
}
