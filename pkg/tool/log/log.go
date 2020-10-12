package log

import (
	"fmt"
	"os"
	"time"
)

//#2 Issue
func WriteLog(data string) {
	file, err := os.OpenFile("test.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	fmt.Fprintln(file, time.Now().Format("2006-01-02T15:04:05+09:00 | ")+data)
}
