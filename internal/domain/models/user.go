package models

type User struct {
	ID       int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Username string `gorm:"unique;not null"`
	PassHash []byte `gorm:"not null"`
	IsAdmin  bool   `gorm:"default:false"`
}
