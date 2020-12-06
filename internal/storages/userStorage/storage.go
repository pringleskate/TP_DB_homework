package userStorage

import (
	//"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/pringleskate/TP_DB_homework/internal/models"
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
)

type Storage interface {
	CreateUser(input models.User) (user models.User, err error)
	GetProfile(input models.UserInput) (user models.User, err error)
	UpdateProfile(nickname models.UserInput, input models.User) (user models.User, err error)
	GetUsers(input models.ForumGetUsers, forumID int) (users []models.User, err error)
//	GetUsers(userIDs []int, conditions models.ForumGetUsers) (users []models.User, err error)
}

type storage struct {
	db *pgx.ConnPool
}

/* constructor */
//func NewStorage(db *pgxpool.Pool) Storage {
func NewStorage(db *pgx.ConnPool) Storage {
	return &storage{
		db: db,
	}
}

//LIMIT - делаем всегда (выставляем максимальное значение int32)
var (
	selectEmpty = "SELECT u.nickname, u.fullname, u.about, u.email FROM forum_users fu JOIN users u ON fu.userID = u.ID WHERE fu.forumID = $1 ORDER BY u.nickname LIMIT $2"
	selectWithSince = "SELECT u.nickname, u.fullname, u.about, u.email FROM forum_users fu JOIN users u ON fu.userID = u.ID WHERE fu.forumID = $1 AND u.nickname > $2 ORDER BY u.nickname LIMIT $3"
	selectWithDesc = "SELECT u.nickname, u.fullname, u.about, u.email FROM forum_users fu JOIN users u ON fu.userID = u.ID WHERE fu.forumID = $1 ORDER BY u.nickname DESC LIMIT $2"
	selectWithSinceDesc =  "SELECT u.nickname, u.fullname, u.about, u.email FROM forum_users fu JOIN users u ON fu.userID = u.ID WHERE fu.forumID = $1 AND AND u.nickname > $2 ORDER BY u.nickname DESC LIMIT $3"
)

func (s *storage) CreateUser(input models.User) (user models.User, err error) {
	//TODO посмотреть, как pgx будет реагировать на null значения, если что сделать default значения в БД
	_, err = s.db.Exec("INSERT INTO users (nickname, email, fullname, about) VALUES ($1, $2, $3, $4)",
						input.Nickname, input.Fullname, input.Email, input.About)

	if pqErr, ok := err.(pgx.PgError); ok {
		switch pqErr.Code {
		case pgerrcode.UniqueViolation:
			return user, models.Error{Code: "409"}
		default:
			return user, models.Error{Code: "500"}
		}
	}

	user.Nickname = input.Nickname
	user.Fullname = input.Fullname
	user.Email = input.Email
	user.About = input.About

	return
}

func (s *storage) GetProfile(input models.UserInput) (user models.User, err error) {
	user.Nickname = input.Nickname
	err = s.db.QueryRow("SELECT fullname, email, about FROM users WHERE nickname = $1", input.Nickname).
				Scan(&user.Fullname, &user.Email, &user.About)

	if err != nil {
		fmt.Println(err)
		if err == pgx.ErrNoRows {
			return user, models.Error{Code: "404"}

		}
		return user, models.Error{Code: "500"}
	}

	return
}

func (s *storage) UpdateProfile(nickname models.UserInput, input models.User) (user models.User, err error) {
	res, err := s.db.Exec("UPDATE users SET nickname = $1, fullname = $2, email = $3, about = $4 WHERE nickname = $5",
						input.Nickname, input.Fullname, input.Email, input.About, nickname.Nickname)

	if pqErr, ok := err.(pgx.PgError); ok {
		switch pqErr.Code {
		case pgerrcode.UniqueViolation:
			return user, models.Error{Code: "409"}
		default:
			return user, models.Error{Code: "500"}
		}
	}

	// если такой пользователь не найден
	if res.RowsAffected() == 0 {
		return user, models.Error{Code: "404"}
	}

	user.Nickname = input.Nickname
	user.Fullname = input.Fullname
	user.Email = input.Email
	user.About = input.About

	return
}

func (s *storage) GetUsers(input models.ForumGetUsers, forumID int) (users []models.User, err error) {
	var rows *pgx.Rows
	if input.Since == "" && !input.Desc {
		rows, err = s.db.Query(selectEmpty, forumID, input.Limit)
	} else if input.Since == "" && input.Desc {
		rows, err = s.db.Query(selectWithDesc, forumID, input.Limit)
	}  else if input.Since != "" && !input.Desc {
		rows, err = s.db.Query(selectWithSince, forumID, input.Since, input.Limit)
	} else if input.Since != "" && input.Desc {
		rows, err = s.db.Query(selectWithSinceDesc, forumID, input.Since, input.Limit)
	}

	if err != nil {
		return users, models.Error{Code: "500"}
	}

	defer rows.Close()

	for rows.Next() {
		user := models.User{}

		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			return users, models.Error{Code: "500"}
		}

		users = append(users, user)
	}

	return
}