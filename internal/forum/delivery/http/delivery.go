package http

import (
	"context"
	"errors"
	pkgForum "github.com/SlavaShagalov/vk-dbms-project/internal/forum"
	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
	pkgErrors "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/errors"
	pkgHTTP "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/http"
	mw "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/middleware"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type delivery struct {
	serv pkgForum.Service
	log  *zap.Logger
}

func RegisterHandlers(router *httprouter.Router, logger *zap.Logger, serv pkgForum.Service) {
	del := delivery{serv, logger}

	router.POST("/api/forum/:slug", mw.AccessLog(mw.HandleError(del.Create, logger), logger))
	router.POST("/api/forum/:slug/:action", mw.AccessLog(mw.HandleError(del.CreateThread, logger), logger))
	//router.POST("/api/forum/:slug/create", mw.AccessLog(mw.HandleError(del.Create, logger), logger))
	//router.POST("/api/forum/create", mw.AccessLog(mw.HandleError(del.Create, logger), logger))
	router.GET("/api/forum/:slug/details", mw.AccessLog(mw.HandleError(del.Get, logger), logger))

	router.GET("/api/forum/:slug/users", mw.AccessLog(mw.HandleError(del.GetForumUsers, logger), logger))
	router.GET("/api/forum/:slug/threads", mw.AccessLog(mw.HandleError(del.GetForumThreads, logger), logger))
}

func (del *delivery) CreateThread(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	action := p.ByName("action")
	if action != "create" {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	slug := p.ByName("slug")

	body, errCreate := pkgHTTP.ReadBody(r, del.log)
	if errCreate != nil {
		return errCreate
	}

	thread := models.Thread{}
	if err := thread.UnmarshalJSON(body); err != nil {
		return pkgErrors.ErrParseJSON
	}

	thread.Forum = slug
	thread, errCreate = del.serv.CreateThread(&thread)

	if errCreate != nil {
		switch {
		case errors.Is(pkgErrors.ErrUserNotFound, errCreate):
			return errCreate
		case errors.Is(pkgErrors.ErrForumNotFound, errCreate):
			return errCreate
		case errors.Is(pkgErrors.ErrThreadAlreadyExists, errCreate):
			w.WriteHeader(http.StatusConflict)
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	data, err := thread.MarshalJSON()
	if err != nil {
		return pkgErrors.ErrInternal
	}

	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return nil
}

func (del *delivery) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	action := p.ByName("slug")
	if action != "create" {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	body, err := pkgHTTP.ReadBody(r, del.log)
	if err != nil {
		return err
	}

	forum := new(models.Forum)
	err = forum.UnmarshalJSON(body)
	if err != nil {
		return pkgErrors.ErrParseJSON
	}

	forum, errCreate := del.serv.Create(context.Background(), forum)
	if errCreate != nil {
		if errCreate != pkgErrors.ErrForumAlreadyExists {
			return errCreate
		}
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	data, err := forum.MarshalJSON()
	if err != nil {
		return pkgErrors.ErrInternal
	}

	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return nil
}

func (del *delivery) Get(w http.ResponseWriter, _ *http.Request, p httprouter.Params) error {
	slug := p.ByName("slug")
	forum, err := del.serv.Get(context.Background(), slug)
	if err != nil {
		return err
	}

	data, err := forum.MarshalJSON()
	if err != nil {
		return pkgErrors.ErrInternal
	}

	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return nil
}

func (del *delivery) GetForumUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	slug := p.ByName("slug")
	var err error

	queryValues := r.URL.Query()
	limit := 100
	strLimit := queryValues.Get("limit")
	if strLimit != "" {
		limit, err = strconv.Atoi(strLimit)
		if err != nil || limit < 0 {
			return pkgErrors.ErrInvalidLimitParam
		}
	}

	since := queryValues.Get("since")

	desc := false
	strDesc := queryValues.Get("desc")
	if strDesc != "" {
		if strDesc == "true" {
			desc = true
		} else if strDesc != "false" {
			return pkgErrors.ErrInvalidDescParam
		}
	}

	users, err := del.serv.GetForumUsers(context.Background(), slug, limit, since, desc)
	if err != nil {
		return err
	}

	data, err := users.MarshalJSON()
	if err != nil {
		return pkgErrors.ErrInternal
	}

	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return nil
}

func (del *delivery) GetForumThreads(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	slug := p.ByName("slug")
	var err error

	queryValues := r.URL.Query()
	limit := 100
	strLimit := queryValues.Get("limit")
	if strLimit != "" {
		limit, err = strconv.Atoi(strLimit)
		if err != nil || limit < 0 {
			return pkgErrors.ErrInvalidLimitParam
		}
	}

	since := queryValues.Get("since")

	desc := false
	strDesc := queryValues.Get("desc")
	if strDesc != "" {
		if strDesc == "true" {
			desc = true
		} else if strDesc != "false" {
			return pkgErrors.ErrInvalidDescParam
		}
	}

	threads, err := del.serv.GetForumThreads(context.Background(), slug, limit, since, desc)
	if err != nil {
		return err
	}

	data, err := threads.MarshalJSON()
	if err != nil {
		return pkgErrors.ErrInternal
	}

	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return nil
}
