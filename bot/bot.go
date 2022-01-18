package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pravinba9495/go-telegram"
	"github.com/pravinba9495/kryptonite/variables"
)

// Run the telegram bot
func Run(botToken string, password string) {
	if botToken != "" {

		go func() {
			for msg := range OutboundChannel {
				if msg != "" && ChatID != "" {
					log.Printf(msg)
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

				if lastMsg != "" {
					if lastMsg == password {
						ChatID = lastChatId
						OutboundChannel <- "🎉 You are now authorized to receive communication through the bot."
						OutboundChannel <- "Your chatId: " + ChatID + ".\n\nRestart the docker container by adding --chatId=" + ChatID + " command line argument to automatically authorize yourself whenever the bot runs."
					} else if strings.ToLower(lastMsg) == "yes" && lastChatId == ChatID {
						if IsWaitingConfirmation {
							ConfirmationChannel <- true
						}
					} else if strings.ToLower(lastMsg) == "no" && lastChatId == ChatID {
						if IsWaitingConfirmation {
							ConfirmationChannel <- false
						}
					} else if strings.ToLower(lastMsg) == "status" && lastChatId == ChatID {
						OutboundChannel <- fmt.Sprintf("Current Status: %s\n%s", variables.CurrentStatus, variables.Verdict)
					} else if strings.ToLower(lastMsg) == "/start" && lastChatId == ChatID {
						OutboundChannel <- "🎉 You are now authorized to receive communication through the bot."
					} else if strings.ToLower(lastMsg) == "manual" && lastChatId == ChatID {
						variables.BotMode = "MANUAL"
						OutboundChannel <- "Kryptonite is now in manual mode."
					} else if strings.ToLower(lastMsg) == "auto" && lastChatId == ChatID {
						variables.BotMode = "AUTO"
						OutboundChannel <- "Kryptonite is now on autopilot."
					} else if strings.ToLower(lastMsg) == "mode" && lastChatId == ChatID {
						OutboundChannel <- fmt.Sprintf("Kryptonite mode: %s", variables.BotMode)
					} else {
						OutboundChannel <- "Command not understood"
					}
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}
