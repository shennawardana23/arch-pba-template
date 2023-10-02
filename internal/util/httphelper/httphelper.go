package httphelper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/schema"
	"github.com/mochammadshenna/arch-pba-template/internal/model/api"
	"github.com/mochammadshenna/arch-pba-template/internal/state"
	"github.com/mochammadshenna/arch-pba-template/internal/util/exceptioncode"
	"github.com/mochammadshenna/arch-pba-template/internal/util/helper"
	"github.com/mochammadshenna/arch-pba-template/internal/util/logger"
)

var Decoder = schema.NewDecoder()

func init() {
	Decoder.RegisterConverter([]string{}, convertStringCommaSeparated)
}

func convertStringCommaSeparated(value string) reflect.Value {
	return reflect.ValueOf(strings.Split(value, ","))
}

func Read(request *http.Request, result interface{}) error {
	err := Decoder.Decode(result, request.URL.Query())

	if err != nil {
		return parseError(err)
	}

	if request.Method == http.MethodPost || request.Method == http.MethodPut || request.Method == http.MethodPatch {
		jsonDecoder := json.NewDecoder(request.Body)
		err = jsonDecoder.Decode(result)
		if err != nil && err != io.EOF {
			logger.Error(request.Context(), err)
			return api.ErrorResponse{
				Code:    exceptioncode.CodeInvalidRequest,
				Message: err.Error(),
			}
		}
	}

	logger.Info(request.Context(), strings.Replace(fmt.Sprintf("request: %+v", result), "\u0026", "", 1))
	return nil
}

func Write(ctx context.Context, writer http.ResponseWriter, data interface{}) {
	response := api.ApiResponse{
		Header: getHeader(writer),
		Data:   data,
	}
	write(ctx, writer, response)
}

func WriteError(ctx context.Context, writer http.ResponseWriter, errorResponse error) {
	writer.WriteHeader(http.StatusBadRequest)
	response := api.ApiResponse{
		Header: getHeader(writer),
		Error:  errorResponse,
	}
	write(ctx, writer, response)
}

func write(ctx context.Context, writer http.ResponseWriter, response api.ApiResponse) {
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	helper.PanicError(err)
}

func getHeader(writer http.ResponseWriter) api.HeaderResponse {
	headerResponse := api.HeaderResponse{
		ServerTimeMs: time.Now().Unix(),
		RequestId:    writer.Header().Get(string(state.HttpHeaders().RequestId)),
	}

	startTimeHeader := writer.Header().Get(string(state.HttpHeaders().StartTime))
	if len(startTimeHeader) > 0 {
		startTime, _ := strconv.ParseInt(startTimeHeader, 10, 64)
		headerResponse.ProcessTimeMs = time.Since(time.Unix(0, startTime)).Milliseconds()
	}

	return headerResponse
}

func parseError(err error) error {
	errors := []api.ErrorValidate{}
	new := err.(schema.MultiError)
	for i, a := range new {
		errors = append(errors, api.ErrorValidate{
			Key:     i,
			Code:    "VALIDATION",
			Message: a.Error(),
		})
	}
	return api.ErrorResponse{
		Code:    exceptioncode.CodeInvalidValidation,
		Message: "validation error",
		Errors:  errors,
	}
}
