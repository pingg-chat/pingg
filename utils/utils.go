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
