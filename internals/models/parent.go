package models

type Parent struct {
	ID     string `gorm:"column:id;not null;uniqueIndex:idx_unique__parent_id"`
	UserID string `gorm:"column:user_id;not null"`
}

func (*Parent) TableName() string {
	return "parent"
}
