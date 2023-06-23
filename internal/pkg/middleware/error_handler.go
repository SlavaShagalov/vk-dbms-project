package middleware

import (
	pkgErrors "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func HandleError(handler func(w http.ResponseWriter, r *http.Request, p httprouter.Params) error, log *zap.Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("after recover: ", zap.Any("error", err))
			}
		}()

		w.Header().Set("Content-Type", "application/json")
		err := handler(w, r, p)
		if err != nil {
			httpCode, exist := pkgErrors.GetHTTPCodeByError(err)
			if !exist {
				err = errors.Wrap(err, "undefined error")
			}

			if httpCode >= 500 {
				log.Error("Internal Server Error", zap.Int("http_code", httpCode), zap.String("error", err.Error()))
			} else {
				log.Info("Response", zap.Int("http_code", httpCode), zap.String("message", err.Error()))
			}

			response := ErrorResponse{Message: err.Error()}
			body, _ := response.MarshalJSON()
			w.WriteHeader(httpCode)
			_, err = w.Write(body)
			if err != nil {
				log.Error(err.Error())
			}
		}
	}
}
