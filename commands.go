package main

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

var loc, _ = time.LoadLocation("Asia/Singapore")

const (
	helpText = "Hi I'm MechaTham! Nice to meet you!" +
		"\n\nTo begin, type one of these commands:" +
		"\n/indicate - To indicate if you are coming home for dinner tonight" +
		"\n/status - To view who is currently coming back for dinner"
	indicateText   = "Are you coming back for dinner tonight?"
	yesText        = "Ok, Indicated coming back for dinner"
	noText         = "Ok, Not coming back for dinner"
	statusText     = "For today, those having dinner: \n"
	notCommandText = "No such command"
	noCommandText  = "Type a Command to begin or /help for more information."
)

func processUpdate(updateStruct *Update, ctx context.Context) interface{} {
	messageEntities := updateStruct.Message_.Entities
	chatId := updateStruct.Message_.Chat.Id

	// Check for valid id
	var famMem FamilyMember
	famKey := datastore.NewKey(ctx, "FamilyMember", "", int64(chatId), nil)
	err := datastore.Get(ctx, famKey, &famMem)
	if err == datastore.ErrNoSuchEntity {
		return makeMessage_NoKeyboard(chatId, "Still in Development")
	}

	// Check for commands
	if messageEntities != nil && messageEntities[0].Type == "bot_command" {
		message := updateStruct.Message_.Text
		switch message {
		case "/help", "/start":
			return makeMessage_NoKeyboard(chatId, helpText)
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
			return makeMessage_Keyboard(chatId, indicateText, replyKeyboardMarkup)
		case "/Yes":
			q := datastore.NewQuery("DinnerStatus").Ancestor(famKey).Order("Date")
			t := q.Run(ctx)
			var ds DinnerStatus
			key, err := t.Next(&ds)
			today := time.Now().In(loc)
			y1, m1, d1 := ds.Date.In(loc).Date()
			y2, m2, d2 := today.Date()
			if err == datastore.Done || (y1 != y2 || m1 != m2 || d1 != d2) {
				ds.Coming = true
				ds.Date = today
				key = datastore.NewIncompleteKey(ctx, "DinnerStatus", famKey)
				datastore.Put(ctx, key, &ds)
			} else {
				ds.Coming = true
				datastore.Put(ctx, key, &ds)
			}

			return makeMessage_NoKeyboard(chatId, yesText)
		case "/No":
			q := datastore.NewQuery("DinnerStatus").Ancestor(famKey).Order("Date")
			t := q.Run(ctx)
			var ds DinnerStatus
			key, err := t.Next(&ds)
			today := time.Now().In(loc)
			y1, m1, d1 := ds.Date.In(loc).Date()
			y2, m2, d2 := today.Date()
			if err == datastore.Done || (y1 != y2 || m1 != m2 || d1 != d2) {
				ds.Coming = false
				ds.Date = today
				key = datastore.NewIncompleteKey(ctx, "DinnerStatus", famKey)
				datastore.Put(ctx, key, &ds)
			} else {
				ds.Coming = false
				datastore.Put(ctx, key, &ds)
			}

			return makeMessage_NoKeyboard(chatId, noText)
		case "/status":
			var statuses string
			q := datastore.NewQuery("FamilyMember").Order("BornYear")
			t := q.Run(ctx)
			for {
				var fm FamilyMember
				key, err := t.Next(&fm)
				if err == datastore.Done {
					break
				}
				q2 := datastore.NewQuery("DinnerStatus").Ancestor(key).Order("Date")
				t2 := q2.Run(ctx)
				var ds DinnerStatus
				_, err2 := t2.Next(&ds)
				today := time.Now().In(loc)
				y1, m1, d1 := ds.Date.In(loc).Date()
				y2, m2, d2 := today.Date()
				if err2 == datastore.Done || (y1 != y2 || m1 != m2 || d1 != d2) {
					statuses += fm.Name + ": " + "Have Not Replied\n"
				} else {
					var coming string
					if ds.Coming {
						coming = "Coming"
					} else {
						coming = "Not Coming"
					}
					statuses += fm.Name + ": " + coming
				}
			}
			return makeMessage_NoKeyboard(chatId, statusText+statuses)
		default:
			return makeMessage_NoKeyboard(chatId, notCommandText)
		}
	} else {
		return makeMessage_NoKeyboard(chatId, "Hi "+famMem.Name+"! "+noCommandText)
	}
}

// Make Messages
const sendMessageMethod = "sendMessage"

func makeMessage_NoKeyboard(chatId int, reply string) SendMessage_NoKeyboard {
	return SendMessage_NoKeyboard{
		Chat_id: chatId,
		Text:    reply,
		Method:  sendMessageMethod,
	}
}

func makeMessage_Keyboard(chatId int, reply string, replyKeyboardMarkup ReplyKeyboardMarkup) SendMessage_Keyboard {
	return SendMessage_Keyboard{
		Chat_id:      chatId,
		Text:         reply,
		Method:       sendMessageMethod,
		Reply_markup: replyKeyboardMarkup,
	}
}
