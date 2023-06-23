package http

import (
	"context"
	"errors"
	pkgErrors "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/errors"
	pkgHTTP "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/http"
	mw "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/middleware"
	pkgUser "github.com/SlavaShagalov/vk-dbms-project/internal/user"
	"go.uber.org/zap"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type delivery struct {
	serv pkgUser.Service
	log  *zap.Logger
}

func RegisterHandlers(router *httprouter.Router, logger *zap.Logger, serv pkgUser.Service) {
	del := delivery{serv, logger}

	router.POST("/api/user/:nickname/create", mw.AccessLog(mw.HandleError(del.Create, logger), logger))
	router.GET("/api/user/:nickname/profile", mw.AccessLog(mw.HandleError(del.GetByNickname, logger), logger))
	router.POST("/api/user/:nickname/profile", mw.AccessLog(mw.HandleError(del.Update, logger), logger))
}

func (del *delivery) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	nickname := p.ByName("nickname")

	body, err := pkgHTTP.ReadBody(r, del.log)
	if err != nil {
		return err
	}

	var request createRequest
	err = request.UnmarshalJSON(body)
	if err != nil {
		// TODO: log error
		return pkgErrors.ErrParseJSON
	}

	params := &pkgUser.CreateParams{
		Nickname: nickname,
		Fullname: request.Fullname,
		About:    request.About,
		Email:    request.Email,
	}
	users, err := del.serv.Create(context.Background(), params)
	if err != nil {
		if errors.Is(err, pkgErrors.ErrUserAlreadyExists) {
			response := newCreateAlreadyExistsResponse(users)
			data, err := response.MarshalJSON()
			if err != nil {
				return pkgErrors.ErrInternal
			}

			w.WriteHeader(http.StatusConflict)
			_, err = w.Write(data)
			if err != nil {
				return pkgErrors.ErrInternal
			}
			return nil
		}
		return err
	}

	response := newCreateResponse(&users[0])
	data, err := response.MarshalJSON()
	if err != nil {
		return pkgErrors.ErrInternal
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return nil
}

func (del *delivery) GetByNickname(w http.ResponseWriter, _ *http.Request, p httprouter.Params) error {
	nickname := p.ByName("nickname")
	user, err := del.serv.GetByNickname(context.Background(), nickname)
	if err != nil {
		return err
	}

	response := newGetResponse(user)
	data, err := response.MarshalJSON()
	if err != nil {
		return pkgErrors.ErrInternal
	}

	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return nil
}

func (del *delivery) Update(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	nickname := p.ByName("nickname")

	body, err := pkgHTTP.ReadBody(r, del.log)
	if err != nil {
		return err
	}

	var request updateRequest
	err = request.UnmarshalJSON(body)
	if err != nil {
		return pkgErrors.ErrParseJSON
	}

	params := &pkgUser.UpdateParams{
		Nickname: nickname,
		Fullname: request.Fullname,
		About:    request.About,
		Email:    request.Email,
	}
	user, err := del.serv.Update(context.Background(), params)
	if err != nil {
		return err
	}

	response := newUpdateResponse(user)
	data, err := response.MarshalJSON()
	if err != nil {
		return pkgErrors.ErrInternal
	}

	_, err = w.Write(data)
	if err != nil {
		return pkgErrors.ErrInternal
	}
	return nil
}
