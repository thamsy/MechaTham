package main

// Json Objects
// - Receive from Bot API
type Update struct {
	Update_id      int     `json:"update_id"`
	Message_       Message `json:"message"`
	Edited_message Message `json:"edited_message"`
}

// - Send to Bot API
type SendMessage struct {
	Chat_id int    `json:"chat_id"`
	Text    string `json:"text"`
	Method  string `json:"method"`
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

	Edit_date int    `json:"edit_date"`
	Text      string `json:"text"`
}
