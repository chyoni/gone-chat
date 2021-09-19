package entity

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	CreatedAt int
	UpdatedAt int
	Rooms     []*Room `gorm:"many2many:user_rooms"`
}

type Room struct {
	ID           uint `gorm:"primaryKey"`
	CreatedAt    int
	UpdatedAt    int
	Participants []*User `gorm:"many2many:user_rooms"`
}
