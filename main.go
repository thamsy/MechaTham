package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/appengine"
)

func init() {
	http.HandleFunc("/" /* + BOTTOKEN*/, rootHandler)
	http.HandleFunc("/init" /* + BOTTOKEN*/, initMem)
	//http.ListenAndServe(":8080", nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Method == "GET" {
		fmt.Fprint(w, "Hello world!")
	} else if r.Method == "POST" {
		// Decoding Received JSON
		var updateStruct Update
		err := json.NewDecoder(r.Body).Decode(&updateStruct)

		//Respond to command
		sendMessageStruct := processUpdate(&updateStruct, ctx)

		// Encoding JSON for sending
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&sendMessageStruct)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func initMem(w http.ResponseWriter, r *http.Request) {
	//Initialize Family Members - Puts into the datastore, telegram_id and member_name
	ctx := appengine.NewContext(r)
	err := initializeMembers(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
