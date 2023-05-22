package models

type User struct {
	ID   string `gorm:"column:id;not null;uniqueIndex:idx_unique__user_id"`
	Name string `gorm:"column:name;not null"`
	DOB  string `gorm:"column:dob;not null"`
}

func (*User) TableName() string {
	return "user"
}
