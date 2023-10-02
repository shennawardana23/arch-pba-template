package exception

import (
	"context"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/mochammadshenna/arch-pba-template/internal/model/api"
	"github.com/mochammadshenna/arch-pba-template/internal/util/exceptioncode"
	"github.com/mochammadshenna/arch-pba-template/internal/util/httphelper"
)

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {
	if isDataNotFoundError(request.Context(), writer, err) {
		return
	}

	if isValidationError(request.Context(), writer, err) {
		return
	}

	if isErrorForeignKeyViolation(request.Context(), writer, err) {
		return
	}

	writeResponse(request.Context(), writer, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", err)
}

func isDataNotFoundError(ctx context.Context, writer http.ResponseWriter, err interface{}) bool {
	exception, ok := err.(exceptioncode.ErrorNotFound)
	if ok {
		writeResponse(ctx, writer, http.StatusNotFound, exceptioncode.CodeDataNotFound, exception.ErrorMessage)
		return true
	}
	return false
}

func isValidationError(ctx context.Context, writer http.ResponseWriter, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	if ok {
		writeResponse(ctx, writer, http.StatusBadRequest, "BAD_REQUEST", exception.Error())
		return true
	}
	return false
}

func isErrorForeignKeyViolation(ctx context.Context, writer http.ResponseWriter, err interface{}) bool {
	exception, ok := err.(exceptioncode.ErrorForeignKeyViolation)
	if ok {
		writeResponse(ctx, writer, http.StatusBadRequest, exceptioncode.CodeInvalidRequest, exception.ErrorMessage)
		return true
	}
	return false
}

func writeResponse(ctx context.Context, writer http.ResponseWriter, httpStatus int, errorCode string, err interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpStatus)

	errorResponse := api.ErrorResponse{
		Code:    errorCode,
		Message: err,
	}

	httphelper.WriteError(ctx, writer, errorResponse)
}
