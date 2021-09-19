package entity

type User struct {
	ID        uint `gorm:"primaryKey"`
	Username  string
	Password  string
	CreatedAt int
	UpdatedAt int
}
