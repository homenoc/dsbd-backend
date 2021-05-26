package slack

func NoticeSlackType(slackType uint) string {
	if slackType == 0 {
		return "追加"
	} else if slackType == 1 {
		return "削除"
	} else if slackType == 1 {
		return "更新"
	} else {
		return ""
	}
}
