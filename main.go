// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kkdai/LineBotTemplate/game"
	"github.com/kkdai/LineBotTemplate/tron"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client
var chash chan string

func main() {
	game.Init()
	initSettlement()

	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
	fmt.Println(123)
}

func initSettlement() {
	chash = make(chan string)
	go func() {
		for {
			hash, err := tron.GetNewBlock()
			if err != nil {
				log.Println(err)
			}
			log.Println("hash:" + hash)
			chash <- hash
			time.Sleep(10 * time.Second)
		}
	}()
}

func readLines(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lines string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines += scanner.Text()
	}
	return lines, scanner.Err()
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {

		if event.Type == linebot.EventTypeJoin {

		}

		if event.Type == linebot.EventTypePostback {
			log.Println(event.Postback.Data)

			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("post back data:"+event.Postback.Data)).Do(); err != nil {
				log.Print(err)
			}
		}

		if event.Type == linebot.EventTypeMessage {

			log.Println("user id:", event.Source.UserID)

			switch message := event.Message.(type) {
			// Handle only on text message
			case *linebot.TextMessage:
				log.Println("TextMessage recevied")
				// GetMessageQuota: Get how many remain free tier push message quota you still have this month. (maximum 500)
				/*quota, err := bot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}*/

				time.Sleep(10 * time.Second)
				file, _ := ioutil.ReadFile("flex.json")
				contenter, _ := linebot.UnmarshalFlexMessageJSON(file)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewFlexMessage("123", contenter)).Do(); err != nil {
					log.Println(err)
				}

				time.Sleep(10 * time.Second)
				// message.ID: Msg unique ID
				// message.Text: Msg text

				/*
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("msg ID:"+message.ID+":"+"Get:"+message.Text+" , \n OK! remain message:"+strconv.FormatInt(quota.Value, 10))).Do(); err != nil {
						log.Print(err)
					}*/

			// Handle only on Sticker message
			case *linebot.StickerMessage:
				log.Println("StickerMessage recevied")
				var kw string
				for _, k := range message.Keywords {
					kw = kw + "," + k
				}

				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewStickerMessage("1", "1")).Do(); err != nil {
					log.Print(err)
				}

				/*
					outStickerResult := fmt.Sprintf("收到貼圖訊息: %s, pkg: %s kw: %s  text: %s", message.StickerID, message.PackageID, kw, message.Text)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(outStickerResult)).Do(); err != nil {
						log.Print(err)
					}*/

			case *linebot.AudioMessage:
			case *linebot.FileMessage:
			case *linebot.FlexMessage:
			case *linebot.ImageMessage:
			case *linebot.LocationMessage:
			case *linebot.ImagemapMessage:
			case *linebot.VideoMessage:
			}

		}
	}
}
