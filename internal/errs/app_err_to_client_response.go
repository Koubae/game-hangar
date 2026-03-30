package errs

import (
	"fmt"
	"net/http"

	"github.com/koubae/game-hangar/pkg/common"
	"go.uber.org/zap"
)

func AppErrToClientResponseWithLog(err error, msg string, logger common.Logger) *common.ClientResponseError {
	var responseError *common.ClientResponseError
	var lvl string

	appErr := AsAppError(err)
	if appErr.IsServerErr() {
		lvl = "error"
		responseError = &common.ClientResponseError{
			HTTPCode: http.StatusInternalServerError,
			Message:  "unexpected error occurred",
		}
	} else {
		lvl = "info"
		responseError = &common.ClientResponseError{
			HTTPCode: http.StatusBadRequest,
			Message:  fmt.Sprintf("%s, error: %s", msg, err.Error()),
		}
	}

	logger.L(lvl, msg, zap.Error(err))
	return responseError
}
