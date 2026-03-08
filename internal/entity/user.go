package entity

import "time"

type User struct {
	ID           uint      `gorm:"column:id;primaryKey"`
	Username     string    `gorm:"column:username"`
	Firstname    string    `gorm:"column:firstname"`
	Lastname     string    `gorm:"column:lastname"`
	Role         string    `gorm:"column:role;default:user"`
	IsVerified   bool      `gorm:"column:is_verified;default:false"`
	Email        string    `gorm:"column:email;uniqueIndex"`
	PasswordHash string    `gorm:"column:password_hash"`
	Books        []Book    `gorm:"foreignKey:UserID"`
	Reviews      []Review  `gorm:"foreignKey:UserID"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (u *User) TableName() string {
	return "users"
}
