package models

type User struct {
	ID           string
	Name         string `json:"name" valid:"required,type(string)"`
	Password     string `json:"password" valid:"required,type(string)"`
	RegisteredAT string
}
