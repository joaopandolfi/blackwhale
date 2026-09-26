package utils

import "fmt"

// WARNING: only use it in test packages
func Equals(a any, b any) bool {
	str1, err := JsonString(a)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	str2, err := JsonString(b)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return str1 == str2
}
