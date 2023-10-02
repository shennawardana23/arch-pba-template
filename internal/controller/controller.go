package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PbaController interface {
	FindAllBrandHotel(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
