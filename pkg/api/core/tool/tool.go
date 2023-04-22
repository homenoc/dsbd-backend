package tool

import "time"

func GetNowStr() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05")
}

func DateToStr(d time.Time) string {
	return d.Format("2006-01-02 15:04:05")
}

func ToIntP(d int) *int {
	return &d
}

func ToUintP(d uint) *uint {
	return &d
}
