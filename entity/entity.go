package entity

type User struct {
	ID        uint    `json:"id,omitempty" gorm:"primaryKey"`
	Username  string  `json:"username,omitempty" gorm:"unique;not null"`
	Password  string  `json:"password,omitempty" gorm:"not null"`
	Alias     string  `json:"alias,omitempty"`
	Avatar    string  `json:"avatar,omitempty"`
	CreatedAt int     `json:"created_at,omitempty"`
	UpdatedAt int     `json:"updated_at,omitempty"`
	Rooms     []*Room `json:"rooms,omitempty" gorm:"many2many:user_rooms"`
}

type Room struct {
	ID           uint    `json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt    int     `json:"created_at,omitempty"`
	UpdatedAt    int     `json:"updated_at,omitempty"`
	Participants []*User `gorm:"many2many:user_rooms"`
}

type UserRooms struct {
	UserID uint `json:"user_id"`
	RoomID uint `json:"room_id"`
}
