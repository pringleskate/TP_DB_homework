package models

import (
	"time"
)

type Error struct {
	Code string
	//Message string
}

func (e Error) Error() string {
	return e.Code
	//panic("implement me")
}

type Forum struct {
	Slug string
	Title string
	User string
	Threads int
	Posts int
}

type ForumCreate struct {
	Slug string
	Title string
	User string
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

type User struct {
	Nickname string
	Fullname string
	Email string
	About string
}

type Thread struct {
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
	Forum   string    `json:"forumService"`
	//ForumID int       `json:"-"`
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

type ThreadUpdate struct {
	ThreadInput
	//SlagOrID string `json:"-"`
	Title    string `json:"title"`   // Заголовок ветки обсуждения.
	Message  string `json:"message"` // Описание ветки обсуждения.
}

type Post struct {
	ThreadInput
	//SlagOrID string `json:"-"`
	ID       int64  `json:"id"`       // Идентификатор данного сообщения.
	Parent   int64  `json:"parent"`   // Идентификатор родительского сообщения (0 - корневое сообщение обсуждения).
	Author   string `json:"author"`   // Автор, написавший данное сообщение.
	Message  string `json:"message"`  // Собственно сообщение форума.
	IsEdited bool   `json:"isEdited"` // Истина, если данное сообщение было изменено.
	Forum    string `json:"forum"`    // Идентификатор форума (slug) данного сообещния.
//	Thread   int32  `json:"thread"`   // Идентификатор ветви (id) обсуждения данного сообещния.
	Created  string `json:"created"`  // Дата создания сообщения на форуме.
}