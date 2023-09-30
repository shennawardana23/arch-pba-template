package state

import "os"

type AppWrapper struct {
	Environment string
}

const (
	appEnv = "APP_ENV"
)

var App AppWrapper

func init() {
	App.Environment = os.Getenv(appEnv) // local | stg | prod
}

// override env
func SetEnv(env string) {
	App.Environment = env
}
