package models

type Student struct {
	ID    string `gorm:"column:id;not null;primaryKey"`
	Name  string `gorm:"column:name;not null"`
	Grade int16  `gorm:"column:grade;not null"`
}

func (*Student) TableName() string {
	return "student"
}
