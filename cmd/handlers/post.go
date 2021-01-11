package handlers

import (
	"encoding/json"
	"github.com/pringleskate/TP_DB_homework/internal/models"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"time"
)

func (h handler) PostsCreate(c *fasthttp.RequestCtx) {
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

	posts := make([]models.Post, 0)

	forum, err := h.Threads.GetForumByThread(&threadInput)
	if err != nil {
		status, respErr, _ := h.ConvertError(err)
		h.WriteResponse(c, status, respErr)
		return
	}

	if len(postsInput) == 0 {
		response, _ := json.Marshal(posts)
		h.WriteResponse(c, fasthttp.StatusCreated, response)
		return
	}

	created := time.Now().Format(time.RFC3339Nano)
	posts, err = h.Posts.CreatePosts(threadInput, forum, created, postsInput)
	if err != nil {
		if err.Error() == "404" {
			status, respErr, _ := h.ConvertError(err)
			h.WriteResponse(c, status, respErr)
			return
		}
		status, respErr, _ := h.ConvertError(err)
		h.WriteResponse(c, status, respErr)
		return
	}

	err = h.Forums.UpdatePostsCount(models.ForumInput{Slug: forum}, len(posts))
	if err != nil {
		status, respErr, _ := h.ConvertError(err)
		h.WriteResponse(c, status, respErr)
		return
	}

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

