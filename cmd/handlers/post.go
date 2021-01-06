package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/pringleskate/TP_DB_homework/internal/models"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"time"
)

/*func (h handler) PostsCreate(c *fasthttp.RequestCtx) {
	postsInput := make([]models.PostCreate, 0)
	threadInput := models.ThreadInput{}
	err := json.Unmarshal(c.PostBody(), &postsInput)
	if err != nil {
		log.Println(err)
		return
	}

	slugOrID := SlagOrID(c)
	threadInput.ThreadID = slugOrID.ThreadID
	threadInput.Slug = slugOrID.Slug

	posts, err := h.Service.CreatePosts(postsInput, threadInput)
	if err != nil {
		status, respErr, _ := h.ConvertError(err)
		h.WriteResponse(c, status, respErr)
		return
	}

	response, _ := json.Marshal(posts)

	h.WriteResponse(c, fasthttp.StatusCreated, response)
	return
}
*/
func (h handler) PostsCreate(c *fasthttp.RequestCtx) {
	postsInput := make([]models.PostCreate, 0)
	threadInput := models.ThreadInput{}
	err := json.Unmarshal(c.PostBody(), &postsInput)
	if err != nil {
		log.Println(err)
		return
	}

	/*if len(postsInput) == 0 {
		response, _ := json.Marshal(postsInput)

		h.WriteResponse(c, fasthttp.StatusCreated, response)
		return
	}*/

	slugOrID := SlagOrID(c)
	threadInput.ThreadID = slugOrID.ThreadID
	threadInput.Slug = slugOrID.Slug

/*	posts, err := h.Service.CreatePosts(postsInput, threadInput)
	if err != nil {
		status, respErr, _ := h.ConvertError(err)
		h.WriteResponse(c, status, respErr)
		return
	}*/

	posts := make([]models.Post, 0)

	forum, err := h.ThreadStorage.GetForumByThread(&threadInput)
	if err != nil {
	/*	response, _ := json.Marshal(posts)
		h.WriteResponse(c, fasthttp.StatusCreated, response)*/
		status, respErr, _ := h.ConvertError(err)
		h.WriteResponse(c, status, respErr)
		return
		//return []models.Post{}, err
	}

	if len(postsInput) == 0 {
		response, _ := json.Marshal(posts)
		h.WriteResponse(c, fasthttp.StatusCreated, response)
		return
		//return []models.Post{}, nil
	}

	created := time.Now()
	for _, postInput := range postsInput {
		post := models.Post{
			ThreadInput: threadInput,
			Parent:      postInput.Parent,
			Author:      postInput.Author,
			Message:     postInput.Message,
			Forum:       forum,
			Created:     created,
		}

		if post.Parent != 0 {
			parentThread, err := h.PostStorage.CheckParentPostThread(post.Parent)
			if err != nil {
				fmt.Println(err)
				status, respErr, _ := h.ConvertError(err)
				h.WriteResponse(c, status, respErr)
				return
				//return []models.Post{}, err
			}

			if parentThread != post.ThreadID  {
				status, respErr, _ := h.ConvertError(models.Error{Code:"409"})
				h.WriteResponse(c, status, respErr)
				return
				//return []models.Post{}, models.Error{Code:"409"}
			}
		}

		output, err := h.PostStorage.CreatePost(post)
		if err != nil {
			status, respErr, _ := h.ConvertError(err)
			h.WriteResponse(c, status, respErr)
			return
			//return []models.Post{}, err
		}

		posts = append(posts, output)

		err = h.ForumStorage.UpdatePostsCount(models.ForumInput{Slug: forum})
		if err != nil {
			status, respErr, _ := h.ConvertError(err)
			h.WriteResponse(c, status, respErr)
			return
			//return []models.Post{}, err
		}
	}

	userID, err := h.UserStorage.GetUserIDByNickname(postsInput[0].Author)
	if err != nil {
		status, respErr, _ := h.ConvertError(err)
		h.WriteResponse(c, status, respErr)
		return
		//return []models.Post{}, err
	}

	forumID, err := h.ForumStorage.GetForumID(models.ForumInput{Slug: forum})
	if err != nil {
		status, respErr, _ := h.ConvertError(err)
		h.WriteResponse(c, status, respErr)
		return
		//return []models.Post{}, err
	}

	err = h.ForumStorage.AddUserToForum(userID, forumID)
	if err != nil && err.Error() != "409" {
		status, respErr, _ := h.ConvertError(err)
		h.WriteResponse(c, status, respErr)
		return
		//return []models.Post{}, err
	}

//	return posts, nil

	response, _ := json.Marshal(posts)

	h.WriteResponse(c, fasthttp.StatusCreated, response)
	return
}

func (h handler) PostGet(c *fasthttp.RequestCtx) {
	id, _ := strconv.Atoi(c.UserValue("id").(string))
	related := c.QueryArgs().Peek("related")
	post, err := h.Service.GetPost(id, string(related))
	if err != nil {
		status, respErr, _ := h.ConvertError(err)
		h.WriteResponse(c, status, respErr)
		return
	}

	response, _ := json.Marshal(post)

	h.WriteResponse(c, fasthttp.StatusOK, response)
	return
}

func (h handler) PostUpdate(c *fasthttp.RequestCtx) {
	postInput := &models.PostUpdate{}
	id, _ := strconv.Atoi(c.UserValue("id").(string))
	postInput.ID = int(id)

	err := postInput.UnmarshalJSON(c.PostBody())
	if err != nil {
		log.Println(err)
		return
	}

	post, err := h.Service.UpdatePost(*postInput)
	if err != nil {
		status, respErr, _ := h.ConvertError(err)
		h.WriteResponse(c, status, respErr)
		return
	}

	response, _ := json.Marshal(post)

	h.WriteResponse(c, fasthttp.StatusOK, response)
	return
}

