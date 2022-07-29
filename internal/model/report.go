package model

type Report struct {
	ID    string `json:"id"`
	URL   string `json:"url" validate:"required,url"`
	Email string `json:"email" validate:"required,email"`
}
