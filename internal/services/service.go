package services

import (
	"fmt"
	"github.com/pringleskate/TP_DB_homework/internal/models"
	"github.com/pringleskate/TP_DB_homework/internal/storages/databaseService"
	"github.com/pringleskate/TP_DB_homework/internal/storages/forumStorage"
	"github.com/pringleskate/TP_DB_homework/internal/storages/postStorage"
	"github.com/pringleskate/TP_DB_homework/internal/storages/threadStorage"
	"github.com/pringleskate/TP_DB_homework/internal/storages/userStorage"
	"github.com/pringleskate/TP_DB_homework/internal/storages/voteStorage"
)

type Service interface {
	CreateForum(input models.ForumCreate) (models.Forum, error)
	GetForum(input models.ForumInput) (models.Forum, error)
	GetForumThreads(input models.ForumGetThreads) ([]models.Thread, error)
	GetForumUsers(input models.ForumGetUsers) ([]models.User, error)

	CreateUser(input models.User) ([]models.User, error)
	GetUser(nickname string) (models.User, error)
	UpdateUser(input models.User) (models.User, error)

	CreateThread(input models.Thread) (models.Thread, error)
	ThreadVote(input models.Vote) (models.Thread, error)
	GetThread(input models.ThreadInput) (models.Thread, error)
	UpdateThread(input models.ThreadUpdate) (models.Thread, error)
	GetThreadPosts(input models.ThreadGetPosts) ([]models.Post, error)

	CreatePosts(input []models.PostCreate, slagOrID string) ([]models.Post, error)
	GetPost(id int, related []string) ([]models.PostFull, error)
	UpdatePost(input models.PostUpdate) (models.Post, error)

	Clear()
	Status() models.Status
}

type service struct {
	forumStorage forumStorage.Storage
	threadStorage threadStorage.Storage
	userStorage userStorage.Storage
	postStorage postStorage.Storage
	voteStorage voteStorage.Storage
	databaseService databaseService.Service
}

func NewService(forumStorage forumStorage.Storage, threadStorage threadStorage.Storage, userStorage userStorage.Storage, postStorage postStorage.Storage, voteStorage voteStorage.Storage, databaseService databaseService.Service) Service {
	return &service{
		forumStorage:  forumStorage,
		threadStorage: threadStorage,
		userStorage:   userStorage,
		postStorage:   postStorage,
		voteStorage:   voteStorage,
		databaseService: databaseService,
	}
}

func (s service) CreateForum(input models.ForumCreate) (models.Forum, error) {
	return s.forumStorage.CreateForum(input)
}

func (s service) GetForum(input models.ForumInput) (models.Forum, error) {
	return s.forumStorage.GetDetails(input)
}

func (s service) GetForumThreads(input models.ForumGetThreads) ([]models.Thread, error) {
	err := s.forumStorage.CheckIfForumExists(models.ForumInput{Slug: input.Slug})
	if err != nil {
		return []models.Thread{}, err
	}
	return s.threadStorage.GetThreadsByForum(input)
}

func (s service) GetForumUsers(input models.ForumGetUsers) ([]models.User, error) {
	forumID, err := s.forumStorage.GetForumID(models.ForumInput{Slug: input.Slug})
	if err != nil {
		return []models.User{}, err
	}

	return s.userStorage.GetUsers(input, forumID)
}

func (s service) CreateUser(input models.User) ([]models.User, error) {
	user, err := s.userStorage.CreateUser(input)
	if err == nil {
		return []models.User{user}, err
	}
	if err.Error() == "409" {
		return []models.User{}, err
		//FIXME достать пользовтелей с тем же имейлом или никнеймом
	}
	return []models.User{}, err
}

func (s service) GetUser(nickname string) (models.User, error) {
	return s.userStorage.GetProfile(nickname)
}

func (s service) UpdateUser(input models.User) (models.User, error) {
	return s.userStorage.UpdateProfile(input)
}

func (s service) CreateThread(input models.Thread) (models.Thread, error) {
	thread, err := s.threadStorage.CreateThread(input)
	if err != nil {
		return thread, err
	}

	err = s.forumStorage.UpdateThreadsCount(models.ForumInput{Slug: input.Slug})
	if err != nil {
		return models.Thread{}, err
	}

	return thread, err
}

func (s service) ThreadVote(input models.Vote) (models.Thread, error) {
	thread, err := s.threadStorage.CheckThreadIfExists(input.Thread)
	if err != nil {
		return models.Thread{}, err
	}

	input.Thread = thread

	return s.voteStorage.UpdateVote(input)
}

func (s service) GetThread(input models.ThreadInput) (models.Thread, error) {
	return s.threadStorage.GetDetails(input)
}

func (s service) UpdateThread(input models.ThreadUpdate) (models.Thread, error) {
	return s.threadStorage.UpdateThread(input)
}

func (s service) GetThreadPosts(input models.ThreadGetPosts) ([]models.Post, error) {
	thread, err := s.threadStorage.CheckThreadIfExists(input.ThreadInput)
	if err != nil {
		return []models.Post{}, err
	}

	input.ThreadInput = thread

	return s.postStorage.GetPostsByThread(input)
}

//TODO М убрать отсюда slagOrID, а сделать отдельный объект ThreadInput
func (s service) CreatePosts(input []models.PostCreate, slagOrID string) ([]models.Post, error) {
	panic("implement me")
	//posts :=
}

//TODO что будет иметься в виду в related? Это будут только названия или как?
//TODO почему возвращается слайс PostFull???????? М
func (s service) GetPost(id int, related []string) ([]models.PostFull, error) {
	post := models.PostFull{}
	err := s.postStorage.GetPostDetails(models.PostInput{ID: id}, &post)
	if err != nil {
		return []models.PostFull{}, err
	}

	return []models.PostFull{}, err
	//TODO проверку на элемент слайса (author и тд)
}

func (s service) UpdatePost(input models.PostUpdate) (models.Post, error) {
	return s.postStorage.UpdatePost(input)
}

func (s service) Clear() {
	err := s.databaseService.Clear()
	if err != nil {
		fmt.Println(err)
	}
}

func (s service) Status() models.Status {
	status, err := s.databaseService.Status()
	if err != nil {
		fmt.Println(err)
	}
	return status
}
