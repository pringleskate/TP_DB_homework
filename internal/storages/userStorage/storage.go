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
}

type storage struct {
	db *pgx.ConnPool
	//db *pgxpool.Pool
}

/* constructor */
//func NewStorage(db *pgxpool.Pool) Storage {
func NewStorage(db *pgx.ConnPool) Storage {
	return &storage{
		db: db,
	}
}

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