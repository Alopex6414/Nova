package app

type User struct {
	UserId      string `json:"userId" yaml:"userId" binding:"required"`
	Username    string `json:"username" yaml:"username" binding:"required"`
	Password    string `json:"password" yaml:"password" binding:"required"`
	PhoneNumber string `json:"phone_number" yaml:"phone_number" binding:"required"`
	Email       string `json:"email" yaml:"email" binding:"omitempty"`
	Address     string `json:"address" yaml:"address" binding:"omitempty"`
	Company     string `json:"company" yaml:"company" binding:"omitempty"`
}

type UserLogin struct {
	UserId   string `json:"userId" yaml:"userId" binding:"required"`
	Username string `json:"username" yaml:"username" binding:"required"`
	Password string `json:"password" yaml:"password" binding:"required"`
}

type ProblemDetails struct {
	Type   string `json:"type" yaml:"type"`
	Title  string `json:"title" yaml:"title"`
	Status int    `json:"status" yaml:"status"`
	Cause  string `json:"cause" yaml:"cause"`
}

type QuestionSingleChoice struct {
	Id             string           `json:"id" yaml:"id" binding:"required"`
	Title          string           `json:"title" yaml:"title" binding:"required"`
	Answers        []QuestionAnswer `json:"answers" yaml:"answers" binding:"required"`
	StandardAnswer QuestionAnswer   `json:"standard_answer" yaml:"standard_answer" binding:"required"`
}

type QuestionMultipleChoice struct {
	Id              string           `json:"id" yaml:"id" binding:"required"`
	Title           string           `json:"title" yaml:"title" binding:"required"`
	Answers         []QuestionAnswer `json:"answers" yaml:"answers" binding:"required"`
	StandardAnswers []QuestionAnswer `json:"standard_answers" yaml:"standard_answers" binding:"required"`
}

type QuestionJudgement struct {
	Id             string `json:"id" yaml:"id" binding:"required"`
	Title          string `json:"title" yaml:"title" binding:"required"`
	Answer         bool   `json:"answer" yaml:"answer"`
	StandardAnswer bool   `json:"standard_answer" yaml:"standard_answer"`
}

type QuestionEssay struct {
	Id             string `json:"id" yaml:"id" binding:"required"`
	Title          string `json:"title" yaml:"title" binding:"required"`
	Answer         string `json:"answer" yaml:"answer" binding:"required"`
	StandardAnswer string `json:"standard_answer" yaml:"standard_answer" binding:"required"`
}

type QuestionTitle struct {
	TitleText string `json:"title_text" yaml:"title_text"`
}

type QuestionAnswer struct {
	AnswerMark string `json:"answerMark" yaml:"answerMark" binding:"required"`
	AnswerText string `json:"answerText" yaml:"answerText" binding:"required"`
}
