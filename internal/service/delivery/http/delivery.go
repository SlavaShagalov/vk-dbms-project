package http

import (
	pkgErrors "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/errors"
	mw "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/middleware"
	pkgService "github.com/SlavaShagalov/vk-dbms-project/internal/service"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"net/http"
)

type delivery struct {
	serv pkgService.Service
	log  *zap.Logger
}

func RegisterHandlers(router *httprouter.Router, log *zap.Logger, serv pkgService.Service) {
	del := delivery{serv, log}

	router.POST("/api/service/clear", mw.AccessLog(mw.HandleError(del.Clear, log), log))
	router.GET("/api/service/status", mw.AccessLog(mw.HandleError(del.GetStatus, log), log))
}

func (del *delivery) Clear(_ http.ResponseWriter, _ *http.Request, _ httprouter.Params) error {
	if err := del.serv.Clear(); err != nil {
		return err
	}
	return nil
}

func (del *delivery) GetStatus(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	status, err := del.serv.GetStatus()
	if err != nil {
		return err
	}

	data, err := status.MarshalJSON()
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
