package main

import (
	_ "github.com/go-sql-driver/mysql"
	config "github.com/mochammadshenna/arch-pba-template/config"
	"github.com/mochammadshenna/arch-pba-template/internal/state"
)

func main() {
	config.Init(state.App.Environment)
}
