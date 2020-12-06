package voteStorage

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/pringleskate/TP_DB_homework/internal/models"
)

type Storage interface {
	UpdateVote(vote models.Vote) (thread models.Thread, err error)
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
	insertVote = "INSERT INTO votes (user_nick, voice, thread) VALUES ($1, $2, $3) ON CONFLICT ON CONSTRAINT unique_vote DO UPDATE SET voice = EXCLUDED.voice;"
	updateThreadVotesUp = "UPDATE threads SET votes = votes + 1 WHERE ID = $1 RETURNING ID, author, created, forum, message, slug, title, votes"
	updateThreadVotesDown = "UPDATE threads SET votes = votes - 1 WHERE ID = $1 RETURNING ID, author, created, forum, message, slug, title, votes"
)

func (s *storage) UpdateVote(vote models.Vote) (thread models.Thread, err error) {
	boolVoice := getBoolVoice(vote)
	tx, err := s.db.Begin()
	if err != nil {
		return thread, models.Error{Code: "500"}
	}

	_, err = tx.Exec("SET LOCAL synchronous_commit TO OFF")
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			fmt.Println(txErr)
			return thread, models.Error{Code: "500"}
		}
		return thread, models.Error{Code: "500"}
	}

	_, err = tx.Exec(insertVote, vote.User, boolVoice, vote.Thread)
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			fmt.Println(txErr)
			return thread, models.Error{Code: "500"}
		}
		if pqErr, ok := err.(pgx.PgError); ok {
			switch pqErr.Code {
			case pgerrcode.ForeignKeyViolation:
				return thread, models.Error{Code: "404"}
			default:
				return thread, models.Error{Code: "500"}
			}
		}
		return thread, models.Error{Code: "500"}
	}

	slug := sql.NullString{}
	if boolVoice {
		err = tx.QueryRow(updateThreadVotesUp, vote.Thread).
			     Scan(&thread.ID, &thread.Author, &thread.Created, &thread.Forum, &thread.Message, slug, &thread.Title, &thread.Votes)
	} else {
		err = tx.QueryRow(updateThreadVotesDown, vote.Thread).
			Scan(&thread.ID, &thread.Author, &thread.Created, &thread.Forum, &thread.Message, slug, &thread.Title, &thread.Votes)
	}
	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			fmt.Println(txErr)
			return thread, models.Error{Code: "500"}
		}
		return thread, models.Error{Code: "500"}
	}

	if slug.Valid {
		thread.Slug = slug.String
	}

	if commitErr := tx.Commit(); commitErr != nil {
		fmt.Println(commitErr)
		return thread, models.Error{Code: "500"}
	}

	return
}

func getBoolVoice(vote models.Vote) bool {
	if vote.Voice == 1 {
		return true
	}
	return false
}

