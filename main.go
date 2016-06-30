package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/" /* + BOTTOKEN*/, rootHandler)
	//http.ListenAndServe(":8080", nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		fmt.Fprint(w, "Hello world!")
	} else if r.Method == "POST" {
		var sendMessageStruct interface{}

		// Decoding Received JSON
		var updateStruct Update
		err := json.NewDecoder(r.Body).Decode(&updateStruct)
		chatId := updateStruct.Message_.Chat.Id

		//Respond to command
		messageEntities := updateStruct.Message_.Entities
		if messageEntities != nil {
			messageEntity := messageEntities[0]
			if messageEntity.Type == "bot_command" {
				message := updateStruct.Message_.Text
				switch message {
				case "/help", "/start":
					sendMessageStruct = makeMessage_NoKeyboard(chatId, "Hi I'm MechaTham! Nice to meet you!\n\nTo begin, type one of these commands:\n/indicate - To indicate if you are coming home for dinner tonight\n/status - To view who is currently coming back for dinner")
				case "/indicate":
					replyKeyboardMarkup := ReplyKeyboardMarkup{
						Keyboard: [][]KeyboardButton{
							[]KeyboardButton{
								KeyboardButton{Text: "/Yes"},
								KeyboardButton{Text: "/No"},
							},
						},
						One_time_keyboard: true,
					}
					sendMessageStruct = makeMessage_Keyboard(chatId, "Are you coming back for dinner tonight?", replyKeyboardMarkup)
				case "/Yes":
					sendMessageStruct = makeMessage_NoKeyboard(chatId, "Ok, Indicated coming back for dinner")
				case "/No":
					sendMessageStruct = makeMessage_NoKeyboard(chatId, "Ok, Not coming back for dinner")
				case "/status":
					sendMessageStruct = makeMessage_NoKeyboard(chatId, "these are the ones coming back for dinner")
				default:
					sendMessageStruct = makeMessage_NoKeyboard(chatId, "No such command")
				}
			}
		} else {
			sendMessageStruct = makeMessage_NoKeyboard(chatId, "Type a Command to begin or /help for more information.")
		}

		// Encoding JSON for sending
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&sendMessageStruct)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func makeMessage_NoKeyboard(chatId int, reply string) SendMessage_NoKeyboard {
	var sendMessageStruct = SendMessage_NoKeyboard{
		Chat_id: chatId,
		Text:    reply,
		Method:  "sendMessage",
	}
	return sendMessageStruct
}

func makeMessage_Keyboard(chatId int, reply string, replyKeyboardMarkup ReplyKeyboardMarkup) SendMessage_Keyboard {
	var sendMessageStruct = SendMessage_Keyboard{
		Chat_id:      chatId,
		Text:         reply,
		Method:       "sendMessage",
		Reply_markup: replyKeyboardMarkup,
	}
	return sendMessageStruct
}
