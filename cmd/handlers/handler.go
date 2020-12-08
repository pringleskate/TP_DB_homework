package handlers

import (
	"errors"
	"github.com/pringleskate/TP_DB_homework/internal/models"
	"github.com/pringleskate/TP_DB_homework/internal/services"
	"github.com/valyala/fasthttp"
	"strconv"
)

type Handler interface {
	ForumCreate(c *fasthttp.RequestCtx)
	ForumGet(c *fasthttp.RequestCtx)
	ForumGetThreads(c *fasthttp.RequestCtx)
	ForumGetUsers(c *fasthttp.RequestCtx)

	ThreadCreate(c *fasthttp.RequestCtx)
	ThreadVote(c *fasthttp.RequestCtx)
	ThreadGet(c *fasthttp.RequestCtx)
	ThreadUpdate(c *fasthttp.RequestCtx)
	ThreadGetPosts(c *fasthttp.RequestCtx)

	PostsCreate(c *fasthttp.RequestCtx)
	PostGet(c *fasthttp.RequestCtx)
	PostUpdate(c *fasthttp.RequestCtx)

	UserCreate(c *fasthttp.RequestCtx)
	UserGet(c *fasthttp.RequestCtx)
	UserUpdate(c *fasthttp.RequestCtx)

	Clear(c *fasthttp.RequestCtx)
	Status(c *fasthttp.RequestCtx)
}

type handler struct {
	Service services.Service
}

func NewHandler(Service services.Service) *handler {
	return &handler{
		Service: Service,
	}
}

func (h handler) WriteResponse(c *fasthttp.RequestCtx, status int, body []byte) {
	c.SetContentType("application/json")
	c.SetStatusCode(status)
	c.Write(body)
}

func (h handler) ConvertError(someError error) (status int, body []byte, err error) {
	Error, ok := someError.(models.Error)
	if !ok {
		return status, body, errors.New("it is not server error")
	}

	body, err = Error.Message.MarshalJSON()
	if err != nil {
		return status, body, err
	}

	status, err = strconv.Atoi(err.(models.Error).Code)
	if err != nil {
		return status, body, err
	}

	return status, body, nil
}

func getBool(k string, args *fasthttp.Args) bool {
	v := args.Peek(k)
	if v != nil && v[0] == 't' {
		return true
	}
	return false
}

func SlagOrID(c *fasthttp.RequestCtx) (output models.ThreadInput) {
	slagOrID := c.UserValue("slug_or_id").(string)

	id, err := strconv.Atoi(slagOrID)
	if err != nil {
		output.Slug = slagOrID
		return output
	}
	output.ID = id
	return output
}