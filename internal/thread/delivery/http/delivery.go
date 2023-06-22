package http

import (
	"errors"
	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
	pkgErrors "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/errors"
	pkgHTTP "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/http"
	mw "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/middleware"
	pkgThread "github.com/SlavaShagalov/vk-dbms-project/internal/thread"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type delivery struct {
	serv pkgThread.Service
	log  *zap.Logger
}

func RegisterHandlers(router *httprouter.Router, log *zap.Logger, serv pkgThread.Service) {
	del := delivery{serv, log}

	router.POST("/api/forum/:slug/create", mw.AccessLog(mw.HandleError(del.CreateThread, log), log))
	router.GET("/api/thread/:slug_or_id/details", mw.AccessLog(mw.HandleError(del.GetThread, log), log))
	router.POST("/api/thread/:slug_or_id/details", mw.AccessLog(mw.HandleError(del.UpdateThread, log), log))

	router.POST("/api/thread/:slug_or_id/create", mw.AccessLog(mw.HandleError(del.CreatePost, log), log))
	router.GET("/api/thread/:slug_or_id/posts", mw.AccessLog(mw.HandleError(del.GetPosts, log), log))

	router.POST("/api/thread/:slug_or_id/vote", mw.AccessLog(mw.HandleError(del.AddVote, log), log))
}

func (del *delivery) CreateThread(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	slug := p.ByName("slug")

	body, errCreate := pkgHTTP.ReadBody(r, del.log)
	if errCreate != nil {
		return errCreate
	}

	thread := models.Thread{}
	if err := thread.UnmarshalJSON(body); err != nil {
		return pkgErrors.ErrParseJson
	}

	thread.Forum = slug
	thread, errCreate = del.serv.CreateThread(&thread)

	if errCreate != nil {
		switch {
		case errors.Is(pkgErrors.ErrUserNotFound, errCreate):
			return errCreate
		case errors.Is(pkgErrors.ErrForumNotFound, errCreate):
			return errCreate
		}
	}

	data, err := thread.MarshalJSON()
	if err != nil {
		return pkgErrors.ErrInternal
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return errCreate
}

func (del *delivery) CreatePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	slugOrId := p.ByName("slug_or_id")

	body, err := pkgHTTP.ReadBody(r, del.log)
	if err != nil {
		return err
	}

	var posts models.PostList
	if err := posts.UnmarshalJSON(body); err != nil {
		return pkgErrors.ErrParseJson
	}

	posts, err = del.serv.CreatePosts(slugOrId, posts)
	if err != nil {
		return err
	}

	data, err := posts.MarshalJSON()
	if err != nil {
		return pkgErrors.ErrInternal
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return err
}

func (del *delivery) GetThread(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	slugOrId := p.ByName("slug_or_id")
	thread, err := del.serv.GetThread(slugOrId)
	if err != nil {
		return err
	}

	data, err := thread.MarshalJSON()
	if err != nil {
		return pkgErrors.ErrInternal
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return err
}

func (del *delivery) UpdateThread(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	slugOrId := p.ByName("slug_or_id")

	body, err := pkgHTTP.ReadBody(r, del.log)
	if err != nil {
		return err
	}

	thread := &models.Thread{}
	if err := thread.UnmarshalJSON(body); err != nil {
		return pkgErrors.ErrParseJson
	}

	updatedThread, err := del.serv.UpdateThread(slugOrId, thread)
	if err != nil {
		return err
	}

	data, err := updatedThread.MarshalJSON()
	if err != nil {
		return pkgErrors.ErrInternal
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return err
}

func (del *delivery) GetPosts(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	slugOrId := p.ByName("slug_or_id")
	queryValues := r.URL.Query()
	limit := 30
	strLimit := queryValues.Get("limit")
	if strLimit != "" {
		limit, err := strconv.Atoi(strLimit)
		if err != nil || limit < 0 {
			return pkgErrors.ErrInvalidLimitParam
		}
	}
	strsince := queryValues.Get("since")
	since := 0
	if strsince != "" {
		since, err := strconv.Atoi(strsince)
		if err != nil || since < 0 {
			return pkgErrors.ErrInvalidLimitParam
		}
	}
	sort := queryValues.Get("sort")
	desc := false
	strDesc := queryValues.Get("desc")
	if strDesc != "" {
		desc, err := strconv.Atoi(strDesc)
		if err != nil || desc < 0 {
			return pkgErrors.ErrInvalidLimitParam
		}
	}

	posts, err := del.serv.GetPosts(slugOrId, limit, since, string(sort), desc)
	if err != nil {
		return err
	}

	data, err := posts.MarshalJSON()
	if err != nil {
		return pkgErrors.ErrInternal
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return nil
}

func (del *delivery) AddVote(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	slugOrId := p.ByName("slug_or_id")

	body, err := pkgHTTP.ReadBody(r, del.log)
	if err != nil {
		return err
	}

	vote := models.Vote{}
	if err := vote.UnmarshalJSON(body); err != nil {
		return pkgErrors.ErrParseJson
	}

	thread, err := del.serv.AddVote(slugOrId, &vote)
	if err != nil {
		return err
	}

	data, err := thread.MarshalJSON()
	if err != nil {
		return pkgErrors.ErrInternal
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return nil
}
