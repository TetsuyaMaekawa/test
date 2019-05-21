package mysql

// RcvData ...
type RcvData struct {
	ID   int    `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

// TableName ...
func (t *RcvData) TableName() string {
	return "mytable"
}
