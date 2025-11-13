package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func Dd(v ...interface{}) {
	for _, val := range v {
		b, _ := json.MarshalIndent(val, "", "  ")
		fmt.Println(string(b))
	}

	os.Exit(0)
}

func Dump(title string, v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(title)
	fmt.Println(string(b))
}
