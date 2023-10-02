migrate_version?=v1.14.0 # version reference https://github.com/amacneil/dbmate/releases/
migrate_platform?=linux
# options:
# - linux
# - macos
# - windows
migrate_arch?=amd64
# options:
# - amd64
# - arm64

env?=local # local | stg | prod

migrate-prepare:
	@rm -rf bin
	@mkdir bin

	curl -L https://github.com/amacneil/dbmate/releases/download/v1.14.0/dbmate-$(migrate_platform)-$(migrate_arch) > ./bin/dbmate
	chmod +x ./bin/dbmate

migrate-new:
	export APP_ENV=$(env) && go run cmd/db/main.go new $(name)

migrate-up:
	export APP_ENV=$(env) && go run cmd/db/main.go up

migrate-down:
	export APP_ENV=$(env) && go run cmd/db/main.go rollback

migrate-status:
	export APP_ENV=$(env) && go run cmd/db/main.go status

build:
	go mod tidy
	go build -o cmd/main cmd/main.go

run:
	export APP_ENV=$(env) && go run cmd/main.go -subscriber=false


connect-db-local:
	./scripts/connect-db-from-local.sh