package postStorage

import (
	"github.com/jackc/pgx"
	"github.com/pringleskate/TP_DB_homework/internal/models"
)

type Storage interface {
	CreatePost(input models.Post) (post models.Post, err error)
	GetPostDetails()
	UpdatePost()
}

type storage struct {
	db *pgx.ConnPool
}

/* constructor */
/*func NewStorage(db *pgx.ConnPool) Storage {
	return &storage{
		db: db,
	}
}
*/
