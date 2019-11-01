package models

import "github.com/jinzhu/gorm"

type Test struct {
	gorm.Model
	Age  int
	Name string
}
