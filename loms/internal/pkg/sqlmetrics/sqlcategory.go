package sqlmetrics

type SQLCategory int

const (
	Select SQLCategory = iota + 1
	Insert
	Update
	Delete
)

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
