package log

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"os"
	"time"
)

//#2 Issue
func WriteLog(user, data string) {
	file, err := os.OpenFile(config.Conf.Log.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	fmt.Fprintln(file, time.Now().Format("2006-01-02 15:04:05")+" , "+user+" , "+data)
}
