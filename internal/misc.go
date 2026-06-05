package utils

import "fmt"

func InArray(s []string, el string) bool {
	for _, v := range s {
		if v == el {
			return true
		}
	}
	return false
}

func CheckErrLite(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
