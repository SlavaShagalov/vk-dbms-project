package http

import (
	"github.com/SlavaShagalov/vk-dbms-project/internal/models"
	pkgErrors "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/errors"
	pkgHTTP "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/http"
	mw "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/middleware"
	pkgPost "github.com/SlavaShagalov/vk-dbms-project/internal/post"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

type delivery struct {
	serv pkgPost.Service
	log  *zap.Logger
}

func RegisterHandlers(router *httprouter.Router, log *zap.Logger, serv pkgPost.Service) {
	del := delivery{serv, log}

	router.GET("/api/post/:id/details", mw.AccessLog(mw.HandleError(del.GetPost, log), log))
	router.POST("/api/post/:id/details", mw.AccessLog(mw.HandleError(del.UpdatePost, log), log))
}

func (del *delivery) GetPost(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	idStr := p.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return pkgErrors.ErrInvalidIDParam
	}

	queryValues := r.URL.Query()
	related := queryValues.Get("related")

	fullPost, err := del.serv.GetPost(id, strings.Split(related, ","))
	if err != nil {
		return err
	}

	data, err := fullPost.MarshalJSON()
	if err != nil {
		return pkgErrors.ErrInternal
	}

	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return nil
}

func (del *delivery) UpdatePost(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	idStr := p.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return pkgErrors.ErrInvalidIDParam
	}

	body, err := pkgHTTP.ReadBody(r, del.log)
	if err != nil {
		return err
	}

	post := models.Post{}
	if err := post.UnmarshalJSON(body); err != nil {
		return pkgErrors.ErrParseJSON
	}

	post.Id = id
	post, err = del.serv.UpdatePost(&post)
	if err != nil {
		return err
	}

	data, err := post.MarshalJSON()
	if err != nil {
		return pkgErrors.ErrInternal
	}

	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return nil
}
