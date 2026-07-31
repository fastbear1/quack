package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func InArray(s []string, el string) bool {
	for _, v := range s {
		if v == el {
			return true
		}
	}
	return false
}

func SortArray[S []E, E any](arr S, less func(a, b int) bool) {
	for i := range len(arr) - 1 {
		for j := i + 1; j > 0 && less(j-1, j); {
			arr[j-1], arr[j] = arr[j], arr[j-1]
			j--
		}
	}
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

func CleanUpDir(dirName string) {
	if _, err := os.Stat(dirName); err == nil {
		_ = os.RemoveAll(dirName)
	}
}
