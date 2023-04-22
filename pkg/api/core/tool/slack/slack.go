package slack

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
	"os"
	"strconv"
	"strings"
)

func NoticeSlackType(slackType uint) string {
	if slackType == 0 {
		return "追加"
	} else if slackType == 1 {
		return "削除"
	} else if slackType == 2 {
		return "更新"
	} else {
		return ""
	}
}

func StartAppSlack() {
	notifySlack := notification.Notification.Slack
	socketMode := socketmode.New(
		notifySlack,
		socketmode.OptionDebug(config.IsDebug),
		socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
	)
	authTest, authTestErr := notifySlack.AuthTest()
	if authTestErr != nil {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN is invalid: %v\n", authTestErr)
		os.Exit(1)
	}
	selfUserId := authTest.UserID

	go func() {
		for event := range socketMode.Events {
			switch event.Type {
			case socketmode.EventTypeEventsAPI:
				socketMode.Ack(*event.Request)

				eventPayload, _ := event.Data.(slackevents.EventsAPIEvent)
				switch eventPayload.Type {
				case slackevents.CallbackEvent:
					switch event := eventPayload.InnerEvent.Data.(type) {
					case *slackevents.MessageEvent:
						if event.User != selfUserId && strings.Contains(event.Text, "こんにちは") {
							_, _, err := notifySlack.PostMessage(
								event.Channel,
								slack.MsgOptionText(
									fmt.Sprintf(":wave: こんにちは <@%v> さん！", event.User),
									false,
								),
							)
							if err != nil {
								log.Printf("Failed to reply: %v", err)
							}
						}
					default:
						socketMode.Debugf("Skipped: %v", event)
					}
				default:
					socketMode.Debugf("unsupported Events API eventPayload received")
				}
			case socketmode.EventTypeSlashCommand:
				cmd, ok := event.Data.(slack.SlashCommand)
				if !ok {
					continue
				}
				inputArray := strings.Split(cmd.Text, " ")
				socketMode.Ack(*event.Request)
				if len(inputArray) < 1 {
					socketMode.PostMessage(cmd.ChannelName, invalidCommand(), slack.MsgOptionReplaceOriginal(cmd.ResponseURL))
					continue
				}
				var msgOption slack.MsgOption = invalidCommand()
				switch inputArray[0] {
				case "echo":
					msgOption = slack.MsgOptionText("hello world", false)
				case "user":
					if len(inputArray) < 2 {
						break
					}
					userID := 0
					userID, _ = strconv.Atoi(inputArray[1])
					if userID == 0 {
						break
					}

					msgOption = getUserInfo(uint(userID))
				case "addr":
					if len(inputArray) < 2 {
						break
					}
					msgOption = getAddrInfo(inputArray[1])
				case "asn":
					if len(inputArray) < 2 {
						break
					}
					asn := 0
					asn, _ = strconv.Atoi(inputArray[1])
					if asn == 0 {
						break
					}
					msgOption = getASNInfo(asn)
				case "help":
					msgOption = getHelpInfo()
				}

				_, _, err := socketMode.PostMessage(cmd.ChannelName, msgOption, slack.MsgOptionReplaceOriginal(cmd.ResponseURL))
				if err != nil {
					log.Println(err)
				}

			default:
				socketMode.Debugf("Skipped: %v", event.Type)
			}
		}
	}()

	socketMode.Run()
}
