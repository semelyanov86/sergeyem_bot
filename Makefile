include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run: run the application
.PHONY: run
run:
	go run bot

## db: connect to the database using mysql
.PHONY: db
db:
	mysql ${BOT_DB_DSN}

## migration-create name=$1: create a new database migration
.PHONY: migration-create
migration-create:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## migrate: apply all up database migrations
.PHONY: migrate
migrate: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database $(addsuffix ${BOT_DB_DSN},mysql://) up

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy and vendor dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	/home/sergey/go/bin/staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# BUILD
# ==================================================================================== #

current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build: build the cmd/api application
.PHONY: build
build:
	@echo 'Building bot...'
	go build -ldflags=${linker_flags} -o=./bin/bot bot
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/bot bot

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

production_host_ip = "serv.sergeyem.ru"

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh easylist@${production_host_ip}

## production/deploy/bot: deploy the api to production
.PHONY: production/deploy/bot
production/deploy/bot:
	rsync -P ./bin/linux_amd64/bot easylist@${production_host_ip}:~/bot
	rsync -rP --delete ./migrations easylist@${production_host_ip}:~
	rsync -P ./remote/production/bot.service easylist@${production_host_ip}:~
	ssh -t easylist@${production_host_ip} '\
		cd ~/bot && make migrate \
		&& sudo mv ~/bot.service /etc/systemd/system/ \
		&& sudo systemctl enable bot \
		&& sudo systemctl restart bot \
		&& sudo service apache2 restart \
	'