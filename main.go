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
		// Decoding Received JSON
		var updateStruct Update
		err := json.NewDecoder(r.Body).Decode(&updateStruct)
		chatId := updateStruct.Message_.Chat.Id

		// Encoding JSON for sending
		var sendMessageStruct = SendMessage{
			Chat_id: chatId,
			Text:    "hi....",
			Method:  "sendMessage",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&sendMessageStruct)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
