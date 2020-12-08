package models

import (
	"time"
)

//easyjson:json
type RespError struct {
	Message string `json:"message"`
}

type Error struct {
	Code string
	Message RespError
}

func (e Error) Error() string {
	return e.Code
}

//easyjson:json
type Forum struct {
	Slug string `json:"slug"`
	Title string `json:"title"`
	User string `json:"user"`
	Threads int `json:"threads"`
	Posts int `json:"posts"`
}

//easyjson:json
type ForumCreate struct {
	Slug string `json:"slug"`
	Title string `json:"title"`
	User string `json:"user"`
}

type ForumInput struct {
	Slug string
}

type ForumGetUsers struct {
	Slug string
	Limit int
	Since string
	Desc bool
}

type ForumGetThreads struct {
	Slug string
	Limit int
	Since string
	Desc bool
}

type UserInput struct {
	Nickname string
}

//easyjson:json
type User struct {
	Nickname string `json:"-"`
	Fullname string `json:"fullname"`
	Email string `json:"email"`
	About string `json:"about"`
}

//easyjson:json
type Thread struct {
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
	Forum   string    `json:"forumService"`
	ID      int       `json:"id"`
	Message string    `json:"message"`
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	Votes   int       `json:"votes"`
}

type ThreadInput struct {
	ID  int
	Slug string
}

//easyjson:json
type ThreadUpdate struct {
	ThreadInput
	Title    string `json:"title"`
	Message  string `json:"message"`
}

type ThreadGetPosts struct {
	ThreadInput
	//Thread int
	Limit int
	Since int
	Sort string
	Desc bool
}

type PostInput struct {
	ID       int  `json:"id"`
}

//easyjson:json
type PostUpdate struct {
	ID       int  `json:"id"`
	Message string `json:"message"`
}
//easyjson:json
type PostCreate struct {
	Parent   int64  `json:"parent"`
	Author   string `json:"author"`
	Message  string `json:"message"`
}
//easyjson:json
type Post struct {
	ThreadInput
	//SlagOrID string `json:"-"`
	ID       int64  `json:"id,omitempty"`       // Идентификатор данного сообщения.
	Parent   int64  `json:"parent,omitempty"`   // Идентификатор родительского сообщения (0 - корневое сообщение обсуждения).
	Author   string `json:"author,omitempty"`   // Автор, написавший данное сообщение.
	Message  string `json:"message,omitempty"`  // Собственно сообщение форума.
	IsEdited bool   `json:"isEdited,omitempty"` // Истина, если данное сообщение было изменено.
	Forum    string `json:"forum,omitempty"`    // Идентификатор форума (slug) данного сообещния.
	//	Thread   int32  `json:"thread"`   // Идентификатор ветви (id) обсуждения данного сообещния.
	Created  string `json:"created,omitempty"`  // Дата создания сообщения на форуме.
}
/*
type Post struct {
	//ThreadInput
	ID       int  `json:"id"`       // Идентификатор данного сообщения.
	Parent   int  `json:"parent"`   // Идентификатор родительского сообщения (0 - корневое сообщение обсуждения).
	Author   string `json:"author"`   // Автор, написавший данное сообщение.
	Message  string `json:"message"`  // Собственно сообщение форума.
	IsEdited bool   `json:"isEdited"` // Истина, если данное сообщение было изменено.
	Forum    string `json:"forum"`    // Идентификатор форума (slug) данного сообещния.
	Thread   int  `json:"thread"`   // Идентификатор ветви (id) обсуждения данного сообещния.
	Created  string `json:"created"`  // Дата создания сообщения на форуме.
}
*/

//type Posts []*Post

//easyjson:json
type PostFull struct {
	Author User `json:"author,omitempty"`
	Forum Forum `json:"forum,omitempty"`
	Post Post `json:"post,omitempty"`
	Thread Thread `json:"thread,omitempty"`
}

//easyjson:json
type Vote struct {
	User string `json:"nickname"`
	Voice int `json:"voice"`
	Thread ThreadInput `json:"_"`
}

//easyjson:json
type Status struct {
	Forum  int32 `json:"forum"`
	Post   int64 `json:"post"`
	Thread int32 `json:"thread"`
	User   int32 `json:"user"`
}