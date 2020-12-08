package services

/*import (
	"github.com/pringleskate/TP_DB_homework/internal/models"
	"github.com/pringleskate/TP_DB_homework/internal/storages/databaseService"
	"github.com/pringleskate/TP_DB_homework/internal/storages/forumStorage"
	"github.com/pringleskate/TP_DB_homework/internal/storages/postStorage"
	"github.com/pringleskate/TP_DB_homework/internal/storages/threadStorage"
	"github.com/pringleskate/TP_DB_homework/internal/storages/userStorage"
	"github.com/pringleskate/TP_DB_homework/internal/storages/voteStorage"
)

//TODO спросить у Маши о заполнении структур в хэндлерах
type Service interface {
	CreateForum(inputForum models.ForumCreate) (forum models.Forum, err error)
	GetForumDetails(inputForum models.ForumInput) (forum models.Forum, err error)
	GetForumThreads(inputForum models.ForumInput) (threads []models.Thread, err error)
	GetForumUsers(inputForum models.ForumInput) (users []models.User, err error)

	//в треде будет threadInput (будем вытаскивать id и по нему вызывать другие функции)
	CreateThread(input models.Thread) (thread models.Thread, err error)
	UpdateThread(input models.ThreadUpdate) (thread models.Thread, err error)
	GetThreadDetails(input models.ThreadInput) (thread models.Thread, err error)
	//проверить сущестсование треда
	GetThreadPosts(input models.ThreadGetPosts) (posts models.Posts, err error)
	VoteForThread(input models.Vote) (thread models.Thread, err error)
	//TODO service/clear
	//TODO service/status
	CreatePosts(input models.Posts) (posts models.Posts, err error)
	UpdatePost(input models.PostUpdate) (post models.Post, err error)
//TODO 	GetPostDetails(input models.PostInput, post *models.PostFull) (err error)

	CreateUser()
	GetUserProfile()
	UpdateUserProfile()
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

func (s *service) CreateForum(inputForum models.ForumCreate) (forum models.Forum, err error) {
	return
}
func (s *service) GetForumDetails(inputForum models.ForumInput) (forum models.Forum, err error) {
	return
}
func (s *service) GetForumThreads(inputForum models.ForumInput) (threads []models.Thread, err error) {
	return
}
func (s *service) GetForumUsers(inputForum models.ForumInput) (users []models.User, err error){
	return
}

//в треде будет threadInput (будем вытаскивать id и по нему вызывать другие функции)
func (s *service) CreateThread(input models.Thread) (thread models.Thread, err error){
	return
}
func (s *service) UpdateThread(input models.ThreadUpdate) (thread models.Thread, err error) {
	return
}
func (s *service) GetThreadDetails(input models.ThreadInput) (thread models.Thread, err error) {
	return
}
//проверить сущестсование треда
func (s *service) GetThreadPosts(input models.ThreadGetPosts) (posts models.Posts, err error){
	return
}
func (s *service) VoteForThread(input models.Vote) (thread models.Thread, err error){
	return
}
//service/clear
//service/status
func (s *service) CreatePosts(input models.Posts) (posts models.Posts, err error){
	return
}
func (s *service) UpdatePost(input models.PostUpdate) (post models.Post, err error){
	return
}
//TODO 	GetPostDetails(input models.PostInput, post *models.PostFull) (err error)

func (s *service) CreateUser(){
	return
}
func (s *service) GetUserProfile(){
	return
}
func (s *service) UpdateUserProfile(){
	return
}*/