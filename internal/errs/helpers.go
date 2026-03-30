package errs

import (
	"fmt"

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
			HTTPCode: appErr.GetDefaultCode(),
			Message:  "unexpected error occurred",
		}
	} else {
		lvl = "info"
		responseError = &common.ClientResponseError{
			HTTPCode: appErr.GetDefaultCode(),
			Message:  fmt.Sprintf("%s, error: %s", msg, err.Error()),
		}
	}

	logger.L(lvl, msg, zap.Error(err))
	return responseError
}

func DTOSchemaValidation(dto any) *AppError {
	if err := common.DTOSchemaValidation(dto); err != nil {
		return Wrap(InvalidPayload, err)
	}
	return nil
}
