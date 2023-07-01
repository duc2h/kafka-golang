package models

type Class struct {
	ID   string `gorm:"column:id;not null;primaryKey`
	Name string `gorm:"column:name;not null"`
}

func (*Class) TableName() string {
	return "class"
}
