package main

// JSON OBJECTS
// - Receive from Bot API
type Update struct {
	Update_id      int     `json:"update_id"`
	Message_       Message `json:"message"`
	Edited_message Message `json:"edited_message"`
}

// - Send to Bot API
type SendMessage_NoKeyboard struct {
	Chat_id    int    `json:"chat_id"`
	Text       string `json:"text"`
	Parse_mode string `json:"parse_mode"`
	Method     string `json:"method"`
}

type SendMessage_Keyboard struct {
	Chat_id      int         `json:"chat_id"`
	Text         string      `json:"text"`
	Parse_mode   string      `json:"parse_mode"`
	Method       string      `json:"method"`
	Reply_markup interface{} `json:"reply_markup"`
}

// -- Other Usual structs
type User struct {
	Id         int    `json:"id"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Username   string `json:"username"`
}

type Chat struct {
	Id         int    `json:"id"`
	Type       string `json:"type"`
	Username   string `json:"username"`
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}

type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Url    string `json:"url"`
	User_  User   `json:"user"`
}

type Message struct {
	Message_id        int  `json:"message_id"`
	From              User `json:"from"`
	Date              int  `json:"date"`
	Chat              Chat `json:"chat"`
	Forward_from      User `json:"forward_from"`
	Forward_from_chat User `json:"forward_from_chat"`
	Forward_date      int  `json:"forward_date"`
	//Reply_to_message  Message         `json:"reply_to_message"`
	Edit_date int             `json:"edit_date"`
	Text      string          `json:"text"`
	Entities  []MessageEntity `json:"entities"`
}

type ReplyKeyboardMarkup struct {
	Keyboard          [][]KeyboardButton `json:"keyboard"`
	Resize_keyboard   bool               `json:"resize_keyboard"`
	One_time_keyboard bool               `json:"one_time_keyboard"`
	Selective         bool               `json:"selective"`
}

type KeyboardButton struct {
	Text             string `json:"text"`
	Request_contact  bool   `json:"request_contact"`
	Request_location bool   `json:"request_location"`
}

type ReplyKeyboardHide struct {
	Hide_keyboard bool `json:"hide_keyboard"`
}
