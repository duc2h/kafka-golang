package models

type Student struct {
	ID     string `gorm:"column:id;not null;uniqueIndex:idx_unique__student_id"`
	UserId string `gorm:"column:user_id;not null"`
	Grad   int16  `gorm:"column:grad;not null"`
}

func (*Student) TableName() string {
	return "student"
}
