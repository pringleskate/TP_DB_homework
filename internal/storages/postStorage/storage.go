package postStorage

import (
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/pringleskate/TP_DB_homework/internal/models"
)

type Storage interface {
	CreatePost(input models.Post) (post models.Post, err error)
	GetPostDetails(input models.PostInput, post *models.PostFull) (err error)
	UpdatePost(input models.PostUpdate) (post models.Post, err error)
	GetPostsByThread(input models.ThreadGetPosts) (posts models.Posts, err error)
}

type storage struct {
	db *pgx.ConnPool
}

func NewStorage(db *pgx.ConnPool) Storage {
	return &storage{
		db: db,
	}
}

//TODO в сервисе проверку на slug or id треда, если что вытащить ID из threadStorage, плюс вытаскиваем forum slug
//TODO service - добавить обновление счетчика в форуме
func (s *storage) CreatePost(input models.Post) (post models.Post, err error) {
	if input.Parent == 0 {
		err = s.db.QueryRow("INSERT INTO posts (author, created, forum, message, parent, thread, path) VALUES ($1,$2,$3,$4,$5,$6, array[(select currval('test_posts_id_seq')::integer)]) RETURNING ID",
			input.Author, input.Created, input.Forum, input.Message, input.Parent, input.Thread).Scan(&post.ID)
	} else {
		err = s.db.QueryRow("INSERT INTO posts (author, created, forum, message, parent, thread, path) VALUES ($1,$2,$3,$4,$5,$6, (SELECT path FROM test_posts WHERE id = $5) || (select currval('test_posts_id_seq')::integer)) RETURNING ID",
			input.Author, input.Created, input.Forum, input.Message, input.Parent, input.Thread).Scan(&post.ID)
	}

	if pqErr, ok := err.(pgx.PgError); ok {
		switch pqErr.Code {
		case pgerrcode.UniqueViolation:
			return post, models.Error{Code: "409"}
		case pgerrcode.NotNullViolation, pgerrcode.ForeignKeyViolation:
			return post, models.Error{Code: "404"}
		default:
			return post, models.Error{Code: "500"}
		}
	}
	post.Author = input.Author
	post.Created = input.Created
	post.Forum = input.Forum
	post.Message = input.Message
	post.Parent = input.Parent
	post.Thread = input.Thread
	post.IsEdited = false
	return
}

//TODO заполнение остальной информации о посте в сервисе
func (s *storage) GetPostDetails(input models.PostInput, post *models.PostFull) (err error) {
	err = s.db.QueryRow("SELECT author, created, forum, message, ID , edited, parent, thread FROM posts WHERE ID = $1", input.ID).
				Scan(&post.Post.Author, &post.Post.Created, &post.Post.Forum, &post.Post.Message, &post.Post.ID, &post.Post.IsEdited, &post.Post.Parent, &post.Post.Thread)
	if err != nil {
		return models.Error{Code: "500"}
	}
	return
}

func (s *storage) UpdatePost(input models.PostUpdate) (post models.Post, err error) {
	err = s.db.QueryRow("UPDATE posts SET message = $1, edited = $2 WHERE ID = $3 RETURNING author, created, forum, message, ID , edited, parent, thread", input.Message, true, input.ID).
				Scan(&post.Author, &post.Created, &post.Forum, &post.Message, &post.ID, &post.IsEdited, &post.Parent, &post.Thread)
	if err != nil {
		if err == pgx.ErrNoRows {
			return post, models.Error{Code: "404"}

		}
		return post, models.Error{Code: "500"}
	}
	return
}

const selectPostsFlatLimitByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent, p.thread, p.forum
	FROM post p
	WHERE p.thread = $1
	ORDER BY p.created, p.id
	LIMIT $2
`

const selectPostsFlatLimitDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent, p.thread, p.forum
	FROM post p
	WHERE p.thread = $1
	ORDER BY p.created DESC, p.id DESC
	LIMIT $2
`
const selectPostsFlatLimitSinceByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent, p.thread, p.forum
	FROM post p
	WHERE p.thread = $1 and p.id > $2
	ORDER BY p.created, p.id
	LIMIT $3
`
const selectPostsFlatLimitSinceDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent, p.thread, p.forum
	FROM post p
	WHERE p.thread = $1 and p.id < $2
	ORDER BY p.created DESC, p.id DESC
	LIMIT $3
`
const selectPostsTreeLimitByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent, p.thread, p.forum
	FROM post p
	WHERE p.thread = $1
	ORDER BY p.path
	LIMIT $2
`
const selectPostsTreeLimitDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent, p.thread, p.forum
	FROM post p
	WHERE p.thread = $1
	ORDER BY path DESC
	LIMIT $2
`
const selectPostsTreeLimitSinceByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent, p.thread, p.forum
	FROM post p
	WHERE p.thread = $1 and (p.path > (SELECT p2.path from post p2 where p2.id = $2))
	ORDER BY p.path
	LIMIT $3
`
const selectPostsTreeLimitSinceDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent, p.thread, p.forum
	FROM post p
	WHERE p.thread = $1 and (p.path < (SELECT p2.path from post p2 where p2.id = $2))
	ORDER BY p.path DESC
	LIMIT $3
`
const selectPostsParentTreeLimitByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent, p.thread, p.forum
	FROM post p
	WHERE p.thread = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		WHERE p2.thread = $2 AND p2.parent = 0
		ORDER BY p2.path
		LIMIT $3
	)
	ORDER BY path
`
const selectPostsParentTreeLimitDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent, p.thread, p.forum
	FROM post p
	WHERE p.thread = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		WHERE p2.parent IS NULL and p2.thread = $2
		ORDER BY p2.path DESC
		LIMIT $3
	)
	ORDER BY p.path[1] DESC, p.path[2:]
`

const selectPostsParentTreeLimitSinceByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent, p.thread, p.forum
	FROM post p
	WHERE p.thread = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		WHERE p2.thread = $2 AND p2.parent = 0 and p2.path[1] > (SELECT p3.path[1] from post p3 where p3.id = $3)
		ORDER BY p2.path
		LIMIT $4
	)
	ORDER BY p.path
`
const selectPostsParentTreeLimitSinceDescByID = `
	SELECT p.id, p.author, p.created, p.edited, p.message, p.parent, p.thread, p.forum
	FROM post p
	WHERE p.thread = $1 and p.path[1] IN (
		SELECT p2.path[1]
		FROM post p2
		WHERE p2.thread = $2 AND p2.parent = 0 and p2.path[1] < (SELECT p3.path[1] from post p3 where p3.id = $3)
		ORDER BY p2.path DESC
		LIMIT $4
	)
	ORDER BY p.path[1] DESC, p.path[2:]
`

func (s *storage) GetPostsByThread(input models.ThreadGetPosts) (posts models.Posts, err error){
	var rows *pgx.Rows

	switch input.Sort {
	case "flat":
		if input.Since > 0 {
			if input.Desc {
				rows, err = s.db.Query(selectPostsFlatLimitSinceDescByID, input.Thread,
					input.Since, input.Limit)
			} else {
				rows, err = s.db.Query(selectPostsFlatLimitSinceByID, input.Thread,
					input.Since, input.Limit)
			}
		} else {
			if input.Desc == true {
				rows, err = s.db.Query(selectPostsFlatLimitDescByID, input.Thread, input.Limit)
			} else {
				rows, err = s.db.Query(selectPostsFlatLimitByID, input.Thread, input.Limit)
			}
		}
	case "tree":
		if input.Since > 0 {
			if input.Desc {
				rows, err = s.db.Query(selectPostsTreeLimitSinceDescByID, input.Thread,
					input.Since, input.Limit)
			} else {
				rows, err = s.db.Query(selectPostsTreeLimitSinceByID, input.Thread,
					input.Since, input.Limit)
			}
		} else {
			if input.Desc {
				rows, err = s.db.Query(selectPostsTreeLimitDescByID, input.Thread, input.Limit)
			} else {
				rows, err = s.db.Query(selectPostsTreeLimitByID, input.Thread, input.Limit)
			}
		}
	case "parent_tree":
		if input.Since > 0 {
			if input.Desc {
				rows, err = s.db.Query(selectPostsParentTreeLimitSinceDescByID, input.Thread, input.Thread,
					input.Since, input.Limit)
			} else {
				rows, err = s.db.Query(selectPostsParentTreeLimitSinceByID, input.Thread, input.Thread,
					input.Since, input.Limit)
			}
		} else {
			if input.Desc {
				rows, err = s.db.Query(selectPostsParentTreeLimitDescByID, input.Thread, input.Thread,
					input.Limit)
			} else {
				rows, err = s.db.Query(selectPostsParentTreeLimitByID, input.Thread, input.Thread,
					input.Limit)
			}
		}
	}

	if err != nil {
		return posts, models.Error{Code: "500"}
	}
	defer rows.Close()

	for rows.Next() {
		post := models.Post{}

		err = rows.Scan(&post.ID, &post.Author, &post.Created, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &post.Forum)
		if err != nil {
			return posts, models.Error{Code: "500"}
		}

		posts = append(posts, &post)
	}

	return 
}
