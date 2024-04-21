package sqlmetrics

// SQLCategory это категория запроса в sql
type SQLCategory int

const (
	// Select обозначает категорию select
	Select SQLCategory = iota + 1
	// Insert обозначает категорию insert
	Insert
	// Update обозначает категорию update
	Update
	// Delete обозначает категорию delete
	Delete
)

// String возвращает текстовое представление cat
func (cat SQLCategory) String() string {
	switch cat {
	case Select:
		return "select"
	case Insert:
		return "insert"
	case Update:
		return "update"
	case Delete:
		return "delete"
	default:
		return "undefined"
	}
}
