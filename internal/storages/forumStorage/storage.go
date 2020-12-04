package forumStorage

import (
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/pringleskate/TP_DB_homework/internal/models"
)

type Storage interface {
	CreateForum(forumSlug models.ForumCreate) (forum models.Forum, err error)
	GetDetails(forumSlug models.ForumInput) (forum models.Forum, err error)
	GetUsers(forumSlug models.ForumGetUsers) (users []models.User, err error)
	GetThreads(forumSlug models.ForumGetThreads) (threads []models.Thread, err error)
}

type storage struct {
	//db *pgxpool.Pool
	db *pgx.ConnPool
}

/* constructor */
func NewStorage(db *pgx.ConnPool) Storage {
	return &storage{
		db: db,
	}
}

func (s *storage) CreateForum(forumSlug models.ForumCreate) (forum models.Forum, err error) {
	//TODO service - проверка пользователя
	_, err = s.db.Exec("INSERT INTO forums (slug, title, user_nick) VALUES ($1, $2, $3) ",
						forumSlug.Slug, forumSlug.Title, forumSlug.User)

	if pqErr, ok := err.(pgx.PgError); ok {
		switch pqErr.Code {
		case pgerrcode.UniqueViolation:
			return forum, models.Error{Code: "409"}
		default:
			return forum, models.Error{Code: "500"}
		}
	}

	forum.User = forumSlug.User
	forum.Title = forumSlug.Title
	forum.Slug = forumSlug.Slug
	forum.Posts = 0
	forum.Threads = 0

	//TODO вынести ошибки в константы?
	return forum, nil
}

func (s *storage) GetDetails(forumSlug models.ForumInput) (forum models.Forum, err error) {
	forum.Slug = forumSlug.Slug
	err = s.db.QueryRow("SELECT title, threads, posts, user_nick FROM forums WHERE slug = ", forumSlug.Slug).
				Scan(&forum.Title, &forum.Threads, &forum.Posts, &forum.User)

	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			return forum, models.Error{Code: "404"}

		}
		return forum, models.Error{Code: "500"}
	}

	return forum, nil
}

//FIXME IMPLEMENT
func (s *storage) GetUsers(forumSlug models.ForumGetUsers) (users []models.User, err error) {

	return []models.User{}, nil
}

//FIXME IMPLEMENT
func (s *storage) GetThreads(forumSlug models.ForumGetThreads) (threads []models.Thread, err error) {
	return []models.Thread{}, nil
}
