package utils

import (
	"encoding/json"
	"fmt"
)

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

func PrettyPrint(i any) string {
	s, _ := json.MarshalIndent(i, "", "   ")
	return string(s)
}
