package v0

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/controller"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/hash"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func SendChatAdmin(data controller.Chat) {
	client := &http.Client{}
	client.Timeout = time.Second * 5

	body, _ := json.Marshal(controller.Chat{
		Err:       data.Err,
		CreatedAt: data.CreatedAt,
		UserID:    data.UserID,
		UserName:  data.UserName,
		GroupID:   data.GroupID,
		Admin:     data.Admin,
		Message:   data.Message,
	})

	//Header部分
	header := http.Header{}
	header.Set("Content-Length", "10000")
	header.Add("Content-Type", "application/json")
	header.Add("TOKEN_1", config.Conf.Controller.Auth.Token1)
	header.Add("TOKEN_2", hash.Generate(config.Conf.Controller.Auth.Token2+config.Conf.Controller.Auth.Token3))

	//リクエストの作成
	req, err := http.NewRequest("POST", "http://"+config.Conf.Controller.User.IP+":"+
		strconv.Itoa(config.Conf.Controller.User.Port)+"/api/v1/controller/chat", bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header = header

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
}

func ReceiveChatAdmin(c *gin.Context) {
	token1 := c.Request.Header.Get("TOKEN_1")
	token2 := c.Request.Header.Get("TOKEN_2")

	if err := auth.ControllerAuthentication(controller.Controller{Token1: token1, Token2: token2}); err != nil {
		log.Println(err)
	}

	var input controller.Chat
	log.Println(c.BindJSON(&input))

	support.Broadcast <- support.WebSocketResult{
		CreatedAt: time.Now(),
		UserID:    input.UserID,
		GroupID:   input.GroupID,
		Admin:     input.Admin,
		Message:   input.Message,
	}
}
