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
//	GetThreads(forumSlug models.ForumGetThreads) (threads []models.Thread, err error)
//	GetUsers(forumSlug models.ForumGetUsers) (userIDs []int, err error)
	UpdateThreadsCount(input models.ForumInput) (err error)
	UpdatePostsCount(input models.ForumInput) (err error)
	AddUserToForum(userID int, forumID int) (err error)
	CheckIfForumExists(input models.ForumInput) (err error)
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

func (s *storage) CreateForum(forumSlug models.ForumCreate) (forum models.Forum, err error) {
	_, err = s.db.Exec("INSERT INTO forums (slug, title, user_nick) VALUES ($1, $2, $3) ",
						forumSlug.Slug, forumSlug.Title, forumSlug.User)

	if pqErr, ok := err.(pgx.PgError); ok {
		switch pqErr.Code {
		case pgerrcode.UniqueViolation:
			return forum, models.Error{Code: "409"}
		case pgerrcode.NotNullViolation, pgerrcode.ForeignKeyViolation:
			return forum, models.Error{Code: "404"}
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

//TODO 2v можно сделать в userstorage один запрос с джоинами
//TODO service - check if forum exists
/*func (s *storage) GetUsers(forumSlug models.ForumGetUsers) (userIDs []int, err error) {
	rows, err := s.db.Query("SELECT userID FROM forum_users FU JOIN forums F ON F.forumID = FU.forumID WHERE F.slug = $1", forumSlug.Slug)
	if err != nil {
		return userIDs, models.Error{Code: "500"}
	}
	defer rows.Close()

	for rows.Next() {
		var userID int

		err = rows.Scan(&userID)
		if err != nil {
			return userIDs, models.Error{Code: "500"}
		}

		userIDs = append(userIDs, userID)
	}

	return
}*/

func (s *storage) UpdateThreadsCount(input models.ForumInput) (err error) {
	_, err = s.db.Exec("UPDATE forums SET threads = threads + 1 WHERE slug = $1", input.Slug)
	if err != nil {
		fmt.Println(err)
		return models.Error{Code: "500"}
	}
	return
}

func (s *storage) UpdatePostsCount(input models.ForumInput) (err error) {
	_, err = s.db.Exec("UPDATE forums SET posts = posts + 1 WHERE slug = $1", input.Slug)
	if err != nil {
		fmt.Println(err)
		return models.Error{Code: "500"}
	}
	return
}

//проверка на приналежность пользователя форуму - если нет, то запись добавится (в сервисе провреять на ошибку не 500)
func (s *storage) AddUserToForum(userID int, forumID int) (err error) {
	_, err = s.db.Exec("INSERT INTO forum_users (forumID, userID) VALUES ($1, $2)", forumID, userID)
	if err != nil {
		if pqErr, ok := err.(pgx.PgError); ok {
			switch pqErr.Code {
			case pgerrcode.UniqueViolation:
				return models.Error{Code: "409"}
			}
		}
		return models.Error{Code: "500"}
	}

	return
}

func (s *storage) CheckIfForumExists(input models.ForumInput) (err error) {
	var ID int
	err = s.db.QueryRow("SELECT ID from forums WHERE slug = $1", input.Slug).Scan(&ID)
	if err != nil {
		if err != pgx.ErrNoRows {
			return models.Error{Code: "404"}
		}
		return models.Error{Code: "500"}
	}

	return
}
