package threadStorage

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/pringleskate/TP_DB_homework/internal/models"
)

type Storage interface {
	CreateThread(input models.Thread) (thread models.Thread, err error)
	GetDetails(input models.ThreadInput) (thread models.Thread, err error)
	UpdateThread(input models.ThreadUpdate) (thread models.Thread, err error)
	GetThreadsByForum(input models.ForumGetThreads) (threads []models.Thread, err error)
	CheckThreadIfExists(input models.ThreadInput) (thread models.ThreadInput, err error)
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
//LIMIT - делаем всегда

	selectThreads = "SELECT id, slug, author, created, forum, title, message, votes FROM threads WHERE forum = $1 ORDER BY created LIMIT $2"
	selectThreadsSince = "SELECT id, slug, author, created, forum, title, message, votes FROM threads WHERE forum = $1 AND created >= $2 ORDER BY created LIMIT $3"
	selectThreadsDesc = "SELECT id, slug, author, created, forum, title, message, votes FROM threads WHERE forum = $1 ORDER BY created DESC LIMIT $2"
	selectThreadsSinceDesc =  "SELECT id, slug, author, created, forum, title, message, votes FROM threads WHERE forum = $1 AND created >= $2 ORDER BY created DESC LIMIT $3"

)

func (s *storage) CreateThread(input models.Thread) (thread models.Thread, err error) {
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

	//TODO TAI сделать отдельную функцию по переприсваиванию полей структур
	thread.Slug = input.Slug
	thread.Votes = input.Votes
	thread.Title = input.Title
	thread.Message = input.Message
	thread.Forum = input.Forum
	thread.Created = input.Created
	thread.Author = input.Author
	//TODO service - добавить обновление счетчика в форуме
	//TODO service - добавить пользователя в форум, если он еще не там
	return
}

func (s *storage) GetDetails(input models.ThreadInput) (thread models.Thread, err error) {
	//TODO переписать на запрос с OR (UpdateThread)
	slug := sql.NullString{}
	if input.Slug == "" {
		err = s.db.QueryRow(selectByID, input.ID).
					Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message, &slug, &thread.Title, &thread.Votes)
	} else {
		err = s.db.QueryRow(selectBySlug, input.Slug).
			Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message, &slug, &thread.Title, &thread.Votes)
	}

	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			return thread, models.Error{Code: "404"}

		}
		return thread, models.Error{Code: "500"}
	}

	if slug.Valid {
		thread.Slug = slug.String
	}

	return
}

func (s *storage) UpdateThread(input models.ThreadUpdate) (thread models.Thread, err error) {
	// при обновлении ветки может меняться только message и title
	//в сервисе надо будет заполнять форум

	//TODO slug - nullsqlstring
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

//TODO перед вызовом этой функции проверять в отдельной функции if forum exists
func (s *storage) GetThreadsByForum(input models.ForumGetThreads) (threads []models.Thread, err error) {
	var rows *pgx.Rows
	if input.Since == "" && !input.Desc {
		rows, err = s.db.Query(selectThreads, input.Slug, input.Limit)
	} else if input.Since == "" && input.Desc {
		rows, err = s.db.Query(selectThreadsDesc,  input.Slug, input.Limit)
	}  else if input.Since != "" && !input.Desc {
		rows, err = s.db.Query(selectThreadsSince,  input.Slug, input.Since, input.Limit)
	} else if input.Since != "" && input.Desc {
		rows, err = s.db.Query(selectThreadsSinceDesc,  input.Slug, input.Since, input.Limit)
	}

	if err != nil {
		return threads, models.Error{Code: "500"}
	}
	defer rows.Close()

	for rows.Next() {
		thread := models.Thread{}
		slug := sql.NullString{}

		err = rows.Scan(&thread.ID, &slug, &thread.Author, &thread.Created, &thread.Forum, &thread.Title, &thread.Message, &thread.Votes)
		if err != nil {
			return threads, models.Error{Code: "500"}
		}

		if slug.Valid {
			thread.Slug = slug.String
		}

		threads = append(threads, thread)
	}

	return
}

func (s storage) CheckThreadIfExists(input models.ThreadInput) (thread models.ThreadInput, err error) {
	if input.Slug == "" {
		err = s.db.QueryRow("SELECT ID from threads WHERE ID = $1", input.ID).Scan(&thread.ID)
	} else {
		err = s.db.QueryRow("SELECT ID from threads WHERE slug = $1", input.Slug).Scan(&thread.ID)
	}

	if err != nil {
		if err == pgx.ErrNoRows {
			return thread, models.Error{Code: "404"}
		}
		return thread, models.Error{Code: "500"}
	}

	return
}
