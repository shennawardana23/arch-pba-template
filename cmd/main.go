package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	config "github.com/mochammadshenna/arch-pba-template/config"
	"github.com/mochammadshenna/arch-pba-template/internal/state"
	"github.com/mochammadshenna/arch-pba-template/internal/util/helper"
	"github.com/mochammadshenna/arch-pba-template/internal/util/logger"
)

func main() {
	config.Init(state.App.Environment)
	logger.Init()

	host := fmt.Sprintf("%s:%d", config.Get().Server.Host, config.Get().Server.Port)
	fmt.Printf("Server running on host:%d \n", config.Get().Server.Port)

	// router := routes.NewRouter(PbaController)

	server := http.Server{
		Addr: host,
		// Handler: router,
	}

	err := server.ListenAndServe()
	helper.PanicError(err)
}
