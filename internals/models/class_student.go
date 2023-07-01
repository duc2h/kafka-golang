package models

// the relationship between student and class table. many-many
type ClassStudent struct {
	ClassID   string `gorm:"class_id;not null;primaryKey"`
	StudentID string `gorm:"class_id;not null;primaryKey"`
}

func (*ClassStudent) TableName() string {
	return "class_student"
}
