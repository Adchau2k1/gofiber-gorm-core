package model

type Contact struct {
	BaseModelSoftDelete
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Content  string `json:"content"`
}
