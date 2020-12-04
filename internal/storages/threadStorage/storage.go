package threadStorage

import (
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/pringleskate/TP_DB_homework/internal/models"
)

type Storage interface {
	CreateThread(input models.Thread) (thread models.Thread, err error)
	GetDetails(input models.ThreadInput) (thread models.Thread, err error)
	UpdateThread(input models.ThreadUpdate) (thread models.Thread, err error)
	GetPosts(input models.ThreadInput) (posts []models.Post, err error)
	//TODO Vote
}

type storage struct {
	db *pgx.ConnPool

}

/* constructor */
func NewStorage(db *pgx.ConnPool) Storage {
	return &storage{
		db: db,
	}
}

var (
	insertWithSlug = "INSERT INTO threads (author, created, forum, message, slug, title, votes) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING ID"
	insertWithoutSlug = "INSERT INTO threads (author, created, forum, message, title, votes) VALUES ($1, $2, $3, $4, $5, $6) RETURNING ID"

	selectBySlug = "SELECT FROM threads (author, created, forum, ID, message, slug, title, votes) WHERE slug = $1"
	selectByID = "SELECT FROM threads (author, created, forum, ID, message, slug, title, votes) WHERE ID = $1"
)

func (s *storage) CreateThread(input models.Thread) (thread models.Thread, err error) {
	//FIXME сделать один запрос с OR
	if input.Slug == "" {
		err = s.db.QueryRow(insertWithoutSlug, input.Author, input.Created, input.Forum, input.Message, input.Title, input.Votes).Scan(&thread.ID)
	} else {
		err = s.db.QueryRow(insertWithSlug, input.Author, input.Created, input.Forum, input.Message, input.Slug, input.Title, input.Votes).Scan(&thread.ID)
	}

	if pqErr, ok := err.(pgx.PgError); ok {
		switch pqErr.Code {
		case pgerrcode.UniqueViolation:
			return thread, models.Error{Code: "409"}
		case pgerrcode.NotNullViolation, pgerrcode.ForeignKeyViolation:
			return thread, models.Error{Code: "404"}
		default:
			return thread, models.Error{Code: "500"}
		}
	}

	//TODO сделать отдельную функцию по переприсваиванию полей структур
	thread.Slug = input.Slug
	thread.Votes = input.Votes
	thread.Title = input.Title
	thread.Message = input.Message
	thread.Forum = input.Forum
	thread.Created = input.Created
	thread.Author = input.Author
//	thread.ForumID = input.ForumID
	//TODO добавить обновление счетчика в форуме
	//TODO добавить пользователя в форум
	return
}

func (s *storage) GetDetails(input models.ThreadInput) (thread models.Thread, err error) {
	if input.Slug == "" {
		err = s.db.QueryRow(selectByID, input.ID).
					Scan(&thread.Author, thread.Created, thread.Forum, thread.ID, thread.Message, thread.Slug, thread.Title, thread.Votes)
	} else {
		err = s.db.QueryRow(selectBySlug, input.Slug).
			Scan(&thread.Author, thread.Created, thread.Forum, thread.ID, thread.Message, thread.Slug, thread.Title, thread.Votes)
	}

	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			return thread, models.Error{Code: "404"}

		}
		return thread, models.Error{Code: "500"}
	}

	return
}

func (s *storage) UpdateThread(input models.ThreadUpdate) (thread models.Thread, err error) {

	// при обновлении ветки может меняться только message и title
	//в сервисе надо будет заполнять форум

	if input.Title != "" && input.Message != "" {
		err = s.db.QueryRow("UPDATE threads SET message = $1, title = $2 WHERE ID = $3 OR slug = $4 " +
								"RETURNING author, created, forum, ID, message, slug, title, votes",
							input.Message, input.Title, input.ID, input.Slug).
					Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)

	} else if input.Title != "" && input.Message == "" {
		err = s.db.QueryRow("UPDATE threads SET title = $1 WHERE ID = $2 OR slug = $3 " +
								"RETURNING author, created, forum, ID, message, slug, title, votes",
								input.Title, input.ID, input.Slug).
					Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)

	} else if input.Title == "" && input.Message != "" {
		err = s.db.QueryRow("UPDATE threads SET message = $1 WHERE ID = $2 OR slug = $3 " +
			"RETURNING author, created, forum, ID, message, slug, title, votes",
			input.Message, input.ID, input.Slug).
			Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)


	} else if input.Title == "" && input.Message == "" {
		return
	}

	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			return thread, models.Error{Code: "404"}

		}
		return thread, models.Error{Code: "500"}
	}

	return
}

func (s *storage) GetPosts(input models.ThreadInput) (posts []models.Post, err error) {
	return
}