package main

import (
	"bytes"
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
	"net/http"
	"time"
)

var loc, _ = time.LoadLocation("Asia/Singapore")

const (
	helpText = "Hi I'm *MechaTham*! Nice to meet you!" +
		"\n\nTo begin, type one of these commands:" +
		"\n/indicate - To indicate if you are coming home for dinner tonight" +
		"\n/status - To view who is currently coming back for dinner" +
		"\n/inform - Sends a notification to all those who haven't replied"
	indicateText   = "Are you coming back for dinner tonight?"
	yesText        = "Ok, Indicated coming back for dinner"
	noText         = "Ok, Not coming back for dinner"
	statusText     = "Having dinner: \n\n"
	notCommandText = "No such command"
	noCommandText  = "Type a Command to begin or /help for more information."
	informText     = "Informed"
	errorText      = "Sorry, an Error has occured"
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
		case "/start":
			//initSchedule(ctx, famMem, chatId)
			return makeMessage_NoKeyboard(chatId, helpText)
		case "/help":
			return makeMessage_NoKeyboard(chatId, helpText)
		case "/inform":
			err := informAll(ctx)
			if err != nil {
				return makeMessage_NoKeyboard(chatId, errorText)
			}
			return makeMessage_NoKeyboard(chatId, informText)
		case "/indicate":
			replyKeyboardMarkup := indicate()
			return makeMessage_Keyboard(chatId, indicateText, replyKeyboardMarkup)
		case "/Yes":
			yesOrNo(ctx, famKey, true)
			return makeMessage_NoKeyboard(chatId, yesText)
		case "/No":
			yesOrNo(ctx, famKey, false)
			return makeMessage_NoKeyboard(chatId, noText)
		case "/status":
			statuses := status(ctx)
			return makeMessage_NoKeyboard(chatId, statusText+statuses)
		default:
			return makeMessage_NoKeyboard(chatId, notCommandText)
		}
	} else {
		return makeMessage_NoKeyboard(chatId, "Hi "+famMem.Name+"! "+noCommandText)
	}
}

// Command functions
/*
func initSchedule(ctx context.Context, famMem FamilyMember, chatId int) {
	now := time.Now().In(loc)
	y, m, d := now.Date()
	target := time.Date(y, m, d, 12, 0, 0, 0, loc)
	if target.Before(now) {
		target = target.AddDate(0, 0, 1)
	}
	duration := target.Sub(now)

	go func() {
		time.Sleep(duration)
		startSchedule(ctx, famMem, chatId)
	}()
}
*/
func informAll(ctx context.Context) error {
	q := datastore.NewQuery("FamilyMember").Order("BornYear")
	t := q.Run(ctx)
	for {
		var fm FamilyMember
		key, err := t.Next(&fm)
		if err == datastore.Done {
			break
		}
		q2 := datastore.NewQuery("DinnerStatus").Ancestor(key).Order("-Date")
		t2 := q2.Run(ctx)
		var ds DinnerStatus
		_, err2 := t2.Next(&ds)
		today := time.Now().In(loc)
		y1, m1, d1 := ds.Date.In(loc).Date()
		y2, m2, d2 := today.Date()
		if err2 == datastore.Done || (y1 != y2 || m1 != m2 || d1 != d2) {
			client := urlfetch.Client(ctx)
			uri := "https://api.telegram.org/bot" + BOTTOKEN + "/sendMessage"
			msgStruct := makeMessage_Keyboard(fm.Id, indicateText, indicate())
			b := new(bytes.Buffer)
			json.NewEncoder(b).Encode(msgStruct)
			req, err := http.NewRequest("POST", uri, b)
			req.Header.Set("Content-Type", "application/json")

			_, err = client.Do(req)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func indicate() ReplyKeyboardMarkup {
	return ReplyKeyboardMarkup{
		Keyboard: [][]KeyboardButton{
			[]KeyboardButton{
				KeyboardButton{Text: "/Yes"},
				KeyboardButton{Text: "/No"},
			},
		},
		Resize_keyboard:   true,
		One_time_keyboard: true,
	}
}

func yesOrNo(ctx context.Context, famKey *datastore.Key, option bool) {
	q := datastore.NewQuery("DinnerStatus").Ancestor(famKey).Order("-Date")
	t := q.Run(ctx)
	var ds DinnerStatus
	key, err := t.Next(&ds)
	today := time.Now().In(loc)
	y1, m1, d1 := ds.Date.In(loc).Date()
	y2, m2, d2 := today.Date()
	if err == datastore.Done || (y1 != y2 || m1 != m2 || d1 != d2) {
		ds.Coming = option
		ds.Date = today
		key = datastore.NewIncompleteKey(ctx, "DinnerStatus", famKey)
		datastore.Put(ctx, key, &ds)
	} else {
		ds.Coming = option
		datastore.Put(ctx, key, &ds)
	}
}

func status(ctx context.Context) string {
	var statuses string
	q := datastore.NewQuery("FamilyMember").Order("BornYear")
	t := q.Run(ctx)
	for {
		var fm FamilyMember
		key, err := t.Next(&fm)
		if err == datastore.Done {
			break
		}
		q2 := datastore.NewQuery("DinnerStatus").Ancestor(key).Order("-Date")
		t2 := q2.Run(ctx)
		var ds DinnerStatus
		_, err2 := t2.Next(&ds)
		today := time.Now().In(loc)
		y1, m1, d1 := ds.Date.In(loc).Date()
		y2, m2, d2 := today.Date()
		if err2 == datastore.Done || (y1 != y2 || m1 != m2 || d1 != d2) {
			statuses += fm.Name + " - " + "*Have Not Replied*\n"
		} else {
			var coming string
			if ds.Coming {
				coming = "Yes"
			} else {
				coming = "No"
			}
			statuses += fm.Name + " - " + coming
		}
	}
	return statuses
}

// Make Messages
const sendMessageMethod = "sendMessage"

func makeMessage_NoKeyboard(chatId int, reply string) SendMessage_NoKeyboard {
	return SendMessage_NoKeyboard{
		Chat_id:    chatId,
		Text:       reply,
		Parse_mode: "Markdown",
		Method:     sendMessageMethod,
	}
}

func makeMessage_Keyboard(chatId int, reply string, replyKeyboardMarkup ReplyKeyboardMarkup) SendMessage_Keyboard {
	return SendMessage_Keyboard{
		Chat_id:      chatId,
		Text:         reply,
		Parse_mode:   "Markdown",
		Method:       sendMessageMethod,
		Reply_markup: replyKeyboardMarkup,
	}
}
