package models

type UserAuth struct {
	Id       uint   `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Name     string `json:"name"`
	Email    string `json:"email"gorm:"unique"`
	Password []byte `json:"-"`
}
