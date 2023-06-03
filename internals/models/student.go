package models

type Student struct {
	ID     string `gorm:"column:id;not null;uniqueIndex:idx_unique__student_id"`
	UserId string `gorm:"column:user_id;not null"`
	Grade  int16  `gorm:"column:grade;not null"`
}

func (*Student) TableName() string {
	return "student"
}
