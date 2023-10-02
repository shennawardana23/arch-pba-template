package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/mochammadshenna/arch-pba-template/internal/controller"
	"github.com/mochammadshenna/arch-pba-template/internal/util/exception"
)

func NewRouter(customerController controller.PbaController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/brand", customerController.FindAllBrandHotel)

	router.PanicHandler = exception.ErrorHandler

	return router
}
