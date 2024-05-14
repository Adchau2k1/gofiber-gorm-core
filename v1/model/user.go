package model

type User struct {
	BaseModelSoftDelete
	Username string `json:"username" gorm:"type:varchar(50);<-:create;unique"`
	Password string `json:"password" gorm:"type:text;size:191"`
	IsBanned *bool  `json:"isBanned" gorm:"default:false"`
	Role     string `json:"role" gorm:"default:Member"`
	Fullname string `json:"fullname" gorm:"type:varchar(100)"`
	Image    string `json:"image"`
}

type UserExport struct {
	BaseModelSoftDelete
	Username string `json:"username"`
	IsBanned bool   `json:"isBanned"`
	Role     string `json:"role"`
	Fullname string `json:"fullname"`
	Image    string `json:"image"`
}

func (UserExport) TableName() string {
	return "users"
}
