### --------------------------------------------------------------------------------------------------------------------
### Variables
### (https://www.gnu.org/software/make/manual/html_node/Using-Variables.html#Using-Variables)
### --------------------------------------------------------------------------------------------------------------------

BUILD_DIR ?= build/

NAME=news-api
REPO=github.com/gregbiv/${NAME}

# Custom local environment file
ifneq ("$(wildcard .env)","")
	include .env
	export $(shell sed 's/=.*//' .env)
endif

SRC_DIRS=cmd pkg

BINARY=news-api
BINARY_SRC=$(REPO)/cmd/${NAME}

GO_LINKER_FLAGS=-ldflags="-s -w"

# RAML configuration
RAML_BUILD_DIR ?= "resources/docs"

# Docker enviroment vars
DOCKER_COMPOSE_EXISTS := $(shell command -v docker-compose 2> /dev/null)
DOCKER_CONTAINER=http-api

# Other config
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# Space separated patterns of packages to skip in list, test, format.
IGNORED_PACKAGES := /vendor/

### --------------------------------------------------------------------------------------------------------------------
### RULES
### (https://www.gnu.org/software/make/manual/html_node/Rule-Introduction.html#Rule-Introduction)
### --------------------------------------------------------------------------------------------------------------------
.PHONY: all clean deps build install

all: clean deps build install

# Installs our project: copies binaries
#-----------------------------------------------------------------------------------------------------------------------
install: build-assets install-bin

install-bin:
	@printf "$(OK_COLOR)==> Installing project$(NO_COLOR)\n"
	go install -v $(BINARY_SRC)

# Building
#-----------------------------------------------------------------------------------------------------------------------
build: build-assets build-bin
build-docs: build-docs-api

build-bin:
	@printf "$(OK_COLOR)==> Building$(NO_COLOR)\n"
	@go build -o ${BUILD_DIR}/${BINARY} ${GO_LINKER_FLAGS} ${BINARY_SRC}

build-assets:
	@printf "$(OK_COLOR)==> Building assets$(NO_COLOR)\n"

	@echo " -> Docs"
	@mkdir -p pkg/assets/docs
	@-rm -f pkg/assets/docs/bindata.go
	@go-bindata -o pkg/assets/docs/bindata.go -pkg docs -prefix resources/docs/ resources/docs/...

	@echo " -> Copying migration files"
	@rm -rf ${BUILD_DIR}/resources
	@mkdir -p ${BUILD_DIR}/resources
	@cp -rv ./resources/migrations/ ${BUILD_DIR}/resources/migrations

build-docs-api:
	@printf "$(OK_COLOR)==> API docs$(NO_COLOR)\n"

	@echo "- Generating"
	@mkdir -p "${RAML_BUILD_DIR}"
	@${call raml2html, -i resources/docs/api.raml -o "${RAML_BUILD_DIR}/api.html"}

# Dependencies
#-----------------------------------------------------------------------------------------------------------------------
deps:
	@git config --global url."https://${GITHUB_TOKEN}@github.com/gregbiv/".insteadOf "https://github.com/gregbiv/"
	@git config --global http.https://gopkg.in.followRedirects true

	@echo "$(OK_COLOR)==> Installing glide dependencies$(NO_COLOR)"
	@glide install

deps-dev:
	@printf "$(OK_COLOR)==> Installing Go-bindata$(NO_COLOR)\n"
	@go get -u github.com/jteeuwen/go-bindata/go-bindata

	@printf "$(OK_COLOR)==> Installing Overalls$(NO_COLOR)\n"
	@go get -u github.com/go-playground/overalls

	@printf "$(OK_COLOR)==> Installing CompileDaemon$(NO_COLOR)\n"
	@go get -u github.com/githubnemo/CompileDaemon

	@printf "$(OK_COLOR)==> Installing Linters$(NO_COLOR)\n"
	@go get -u golang.org/x/tools/cmd/goimports
	@go get -u github.com/golang/lint/golint

	@printf "$(OK_COLOR)==> Installing TOC generator$(NO_COLOR)\n"
	@go get -u github.com/nochso/tocenize/cmd/tocenize

# Migrations
#-----------------------------------------------------------------------------------------------------------------------
migrations-dev:
	@printf "$(OK_COLOR)==> Running migrations$(NO_COLOR)\n"
	@docker-compose exec http-api news-api migrate

# Testing
#-----------------------------------------------------------------------------------------------------------------------
test: test-unit

test-unit:
	@printf "$(OK_COLOR)==> Running tests$(NO_COLOR)\n"
	@go test -race -cover -covermode=atomic $(shell go list ./... | grep -v /vendor/)
	@overalls -project=github.com/gregbiv/news-api -ignore="vendor,.glide" -covermode=count
	@curl -s https://codecov.io/bash > codecov.sh
	@bash ./codecov.sh -t ${CODECOV_TOKEN} -f overalls.coverprofile
	@rm codecov.sh

test-dev:
	@printf "$(OK_COLOR)==> Running tests$(NO_COLOR)\n"
	@go test -race -cover -covermode=atomic $(shell go list ./... | grep -v /vendor/)

# Lint
#-----------------------------------------------------------------------------------------------------------------------
lint:
	@echo "$(OK_COLOR)==> Linting... $(NO_COLOR)"
	@golint $(SRC_DIRS)
	@goimports -l -w $(SRC_DIRS)

# Development
#-----------------------------------------------------------------------------------------------------------------------
dev-up:
	@printf "$(OK_COLOR)==> Starting containers$(NO_COLOR)\n"
	@docker-compose up -d

dev-ssh:
	@docker-compose exec "${DOCKER_CONTAINER}" sh

dev-stop:
	@printf "$(OK_COLOR)==> Stopping containers$(NO_COLOR)\n"
	@docker-compose stop

dev-migrate:
	@printf "$(OK_COLOR)==> Executing dev DB migrations$(NO_COLOR)\n"
	@docker-compose exec http-api news-api migrate

# House keeping
#-----------------------------------------------------------------------------------------------------------------------
clean:
	@printf "$(OK_COLOR)==> Cleaning project$(NO_COLOR)\n"
	if [ -d ${BUILD_DIR} ] ; then rm -rf ${BUILD_DIR} ; fi

### --------------------------------------------------------------------------------------------------------------------
### Functions
### --------------------------------------------------------------------------------------------------------------------

# Environment Helpers
#-----------------------------------------------------------------------------------------------------------------------
define protoc-docs
	@protoc $(1)
endef
define raml2html
	@raml2html $(1)
endef

ifdef DOCKER_COMPOSE_EXISTS
define raml2html
	@docker run --rm -v "$$(pwd):/data" "letsdeal/raml2html:6.2" $(1)
endef
endif
