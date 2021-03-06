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

// CONSTANTS
var loc, _ = time.LoadLocation("Asia/Singapore")

const (
	helpText = "Hi I'm *MechaTham*! Nice to meet you!" +
		"\n\nTo begin, type one of these commands:" +
		"\n/indicate - To indicate if you are coming home for dinner tonight" +
		"\n/status - To view who is currently coming back for dinner" +
		"\n/inform - Sends a notification to all those who haven't replied" +
		"\n/remark - Add a Remark to your Dinner status" +
		"\n/cancel - Cancel current operation"
	indicateText      = "Are you coming back for dinner tonight?"
	yesText           = "Ok, Coming Back. Use /remark to add a remark."
	noText            = "Ok, Not Coming Back. Use /remark to add a remark."
	startRemarkText   = "Please type and add your remark, it will be displayed in /status"
	indicateFirstText = "Please indicate if you're coming back for dinner first"
	savedRemarkText   = "Remark saved! View with /status"
	statusText        = "Having dinner: \n\n"
	notCommandText    = "No such command"
	noCommandText     = "Type a Command to begin or /help for more information."
	informText        = "Informed"
	cancelText        = "Current Command cancelled. Type /help for all commands."
	errorText         = "Sorry, an Error has occured"
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
	memPrevComm := famMem.PrevCommand
	message := updateStruct.Message_.Text

	// Check for commands
	if message == "/cancel" {
		setCurrCommand(ctx, famKey, &famMem, "")
		return makeMessage_Keyboard(chatId, cancelText, ReplyKeyboardHide{true})
	} else if memPrevComm != "" {
		switch memPrevComm {
		case "/remark":
			success := saveRemark(ctx, famKey, message)
			if !success { // Defensive coding
				setCurrCommand(ctx, famKey, &famMem, "")
				return makeMessage_Keyboard(chatId, indicateFirstText, indicate())
			} else {
				setCurrCommand(ctx, famKey, &famMem, "")
				return makeMessage_NoKeyboard(chatId, savedRemarkText)
			}
		default:
			return makeMessage_NoKeyboard(chatId, errorText)
		}
	} else if messageEntities != nil && messageEntities[0].Type == "bot_command" {
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
			return makeMessage_Keyboard(chatId, yesText, ReplyKeyboardHide{true})
		case "/No":
			yesOrNo(ctx, famKey, false)
			return makeMessage_Keyboard(chatId, noText, ReplyKeyboardHide{true})
		case "/remark":
			success := saveRemark(ctx, famKey, "")
			if success {
				setCurrCommand(ctx, famKey, &famMem, message)
				return makeMessage_NoKeyboard(chatId, startRemarkText)
			} else {
				return makeMessage_Keyboard(chatId, indicateFirstText, indicate())
			}
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

// COMMAND FUNCTIONS
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

// Sends a message to all family members who hasn't indicated dinner status
func informAll(ctx context.Context) error {
	q := datastore.NewQuery("FamilyMember").Order("BornYear")
	t := q.Run(ctx)
	for {
		var fm FamilyMember
		key, err := t.Next(&fm)
		if err == datastore.Done {
			break
		} /*else if fm.DisableNotifTil.After(time.Now().In(loc)) {
			continue
		}*/
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

// Creates a Keyboard with /Yes and /No
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

// Stores dinner status
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

// Stores Current command
func setCurrCommand(ctx context.Context, famKey *datastore.Key, famMem *FamilyMember, message string) {
	famMem.PrevCommand = message
	datastore.Put(ctx, famKey, famMem)
}

// Save Remark
func saveRemark(ctx context.Context, famKey *datastore.Key, remark string) bool {
	q := datastore.NewQuery("DinnerStatus").Ancestor(famKey).Order("-Date")
	t := q.Run(ctx)
	var ds DinnerStatus
	key, err := t.Next(&ds)
	today := time.Now().In(loc)
	y1, m1, d1 := ds.Date.In(loc).Date()
	y2, m2, d2 := today.Date()
	if err == datastore.Done || (y1 != y2 || m1 != m2 || d1 != d2) {
		return false
	} else {
		ds.Remark = remark
		datastore.Put(ctx, key, &ds)
		return true
	}

}

// Displays dinner status of all family members
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

		// Determine status
		var coming string
		remark := ""
		/*if fm.DisableNotifTil.After(today) {
			coming = "Disabled"
		} else*/if err2 == datastore.Done || (y1 != y2 || m1 != m2 || d1 != d2) {
			coming = "*Have Not Replied*"
		} else {
			if ds.Coming {
				coming = "Yes"
			} else {
				coming = "No"
			}

			if ds.Remark != "" {
				remark = "\n_" + ds.Remark + "_"
			}
		}
		statuses += fm.Name + " - " + coming + remark + "\n"
	}
	return statuses
}

/*func disableNotif(famMem *FamilyMember, famKey *datastore.Key) {

}*/

// MESSAGE OUTLINES
const (
	sendMessageMethod = "sendMessage"
	parseModeMarkdown = "Markdown"
)

func makeMessage_NoKeyboard(chatId int, reply string) SendMessage_NoKeyboard {
	return SendMessage_NoKeyboard{
		Chat_id:    chatId,
		Text:       reply,
		Parse_mode: parseModeMarkdown,
		Method:     sendMessageMethod,
	}
}

func makeMessage_Keyboard(chatId int, reply string, replyKeyboardMarkup interface{}) SendMessage_Keyboard {
	return SendMessage_Keyboard{
		Chat_id:      chatId,
		Text:         reply,
		Parse_mode:   parseModeMarkdown,
		Method:       sendMessageMethod,
		Reply_markup: replyKeyboardMarkup,
	}
}
