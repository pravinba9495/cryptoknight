package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/pravinba9495/go-telegram"
)

// Run the telegram bot
func Run(botToken string, password string) {
	if botToken != "" {

		go func() {
			for msg := range ErrorChannel {
				if msg.Error() != "" && ChatID != "" {
					_, err := telegram.SendMessage(botToken, ChatID, msg.Error())
					if err != nil {
						log.Println(err)
					}
				}
			}
		}()

		go func() {
			for msg := range OutboundChannel {
				if msg != "" && ChatID != "" {
					_, err := telegram.SendMessage(botToken, ChatID, msg)
					if err != nil {
						log.Println(err)
					}
				}
			}
		}()

		offset := 0
		firstRun := false
		for {
			updates, err := telegram.GetUpdates(botToken, fmt.Sprint(offset))
			if err != nil {
				log.Println(err)
			} else {
				var lastMsg, lastChatId string
				for _, update := range *updates {
					lastMsg = update.Message.Text
					lastChatId = fmt.Sprint(update.Message.Chat.ID)
					offset = int(update.UpdateID) + 1
				}
				if !firstRun {
					firstRun = true
					continue
				}

				if lastMsg == password {
					ChatID = lastChatId
					OutboundChannel <- "🎉 You are now authorized to receive communication through the bot."
					OutboundChannel <- "Your chatId: " + ChatID + ".\n\nRestart the docker container by adding --chatId=" + ChatID + " command line argument to automatically authorize yourself whenever the bot runs."
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}