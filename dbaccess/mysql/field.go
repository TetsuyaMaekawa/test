package mysql

// Mydb DBのフィールドを保持
type Mytable struct {
	ID   int
	Name string
}

func (m *Mytable) TableName() string {
	return "mytable"
}
