package http

import (
	"context"
	pkgForum "github.com/SlavaShagalov/vk-dbms-project/internal/forum"
	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
	pkgErrors "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/errors"
	pkgHTTP "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/http"
	mw "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/middleware"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type delivery struct {
	serv pkgForum.Service
	log  *zap.Logger
}

func RegisterHandlers(router *httprouter.Router, logger *zap.Logger, serv pkgForum.Service) {
	del := delivery{serv, logger}

	router.POST("/api/forum/create", mw.AccessLog(mw.HandleError(del.Create, logger), logger))
	router.GET("/api/forum/:slug/details", mw.AccessLog(mw.HandleError(del.Get, logger), logger))
	router.GET("/api/forum/:slug/users", mw.AccessLog(mw.HandleError(del.GetForumUsers, logger), logger))
	router.GET("/api/forum/:slug/threads", mw.AccessLog(mw.HandleError(del.GetForumThreads, logger), logger))
}

func (del *delivery) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	body, err := pkgHTTP.ReadBody(r, del.log)
	if err != nil {
		return err
	}

	forum := new(models.Forum)
	err = forum.UnmarshalJSON(body)
	if err != nil {
		// TODO: log error
		return pkgErrors.ErrParseJson
	}

	forum, errCreate := del.serv.Create(context.Background(), forum)
	if errCreate != nil && errCreate != pkgErrors.ErrForumAlreadyExists {
		return errCreate
	}

	data, err := forum.MarshalJSON()
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

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return nil
}

func (del *delivery) GetForumUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	slug := p.ByName("slug")

	queryValues := r.URL.Query()
	limit := 30
	strLimit := queryValues.Get("limit")
	if strLimit != "" {
		limit, err := strconv.Atoi(strLimit)
		if err != nil || limit < 0 {
			return pkgErrors.ErrInvalidLimitParam
		}
	}
	since := queryValues.Get("since")
	desc := false
	strDesc := queryValues.Get("desc")
	if strDesc != "" {
		desc, err := strconv.Atoi(strDesc)
		if err != nil || desc < 0 {
			return pkgErrors.ErrInvalidLimitParam
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

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return nil
}

func (del *delivery) GetForumThreads(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	slug := p.ByName("slug")

	queryValues := r.URL.Query()
	limit := 30
	strLimit := queryValues.Get("limit")
	if strLimit != "" {
		limit, err := strconv.Atoi(strLimit)
		if err != nil || limit < 0 {
			return pkgErrors.ErrInvalidLimitParam
		}
	}
	since := queryValues.Get("since")

	desc := false
	strDesc := queryValues.Get("desc")
	if strDesc != "" {
		desc, err := strconv.Atoi(strDesc)
		if err != nil || desc < 0 {
			return pkgErrors.ErrInvalidLimitParam
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

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return nil
}
