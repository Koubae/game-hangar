package errs

import (
	"fmt"
	"net/http"

	"github.com/koubae/game-hangar/pkg/common"
	"github.com/koubae/game-hangar/pkg/web"
	"go.uber.org/zap"
)

func AppErrToClientResponseWithLog(w http.ResponseWriter, err error, msg string, logger common.Logger) {
	var lvl string

	if AsAppError(err).IsServerErr() {
		lvl = "error"
	} else {
		lvl = "info"
	}
	logger.L(lvl, msg, zap.Error(err))

	AppErrToClientResponse(w, err, msg)
}

func AppErrToClientResponse(w http.ResponseWriter, err error, msg string) {
	var response common.ClientResponseError

	appErr := AsAppError(err)
	if appErr.IsServerErr() {
		response = common.ClientResponseError{
			HTTPCode: appErr.GetDefaultCode(),
			Message:  "unexpected error occurred",
		}
	} else {
		response = common.ClientResponseError{
			HTTPCode: appErr.GetDefaultCode(),
			Message:  fmt.Sprintf("%s%s", msg, err.Error()),
		}
	}

	web.WriteBusinessErrorResponse(w, &response)
}

func DTOSchemaValidation(dto any) *AppError {
	if err := common.DTOSchemaValidation(dto); err != nil {
		return Wrap(PayloadValidation, err)
	}
	return nil
}
