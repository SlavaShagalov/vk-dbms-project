package http

import (
	"github.com/SlavaShagalov/vk-dbms-project/internal/pkg/constants"
	pkgErrors "github.com/SlavaShagalov/vk-dbms-project/internal/pkg/errors"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func ReadBody(r *http.Request, log *zap.Logger) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(constants.FailedReadRequestBody, zap.Error(err))
		return nil, pkgErrors.ErrReadBody
	}

	err = r.Body.Close()
	if err != nil {
		log.Error(constants.FailedCloseRequestBody, zap.Error(err))
	}

	return body, nil
}
