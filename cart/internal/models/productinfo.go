package models

type ProductInfo struct {
	Name  string
	Price uint32
}

func IsStringValidName(str string) bool {
	return len(str) > 0
}
