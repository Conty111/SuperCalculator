PROJECT_PKG = gitlab.sch.ocrv.com.rzd/blockchain/platform/gobase
BUILD_DIR = build
VERSION ?=$(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE ?= $(shell date +%FT%T%z)
# remove debug info from the binary & make it smaller
LDFLAGS += -s -w
#LDFLAGS += -linkmode external -w -extldflags "-static"
# inject build info
LDFLAGS += -X ${PROJECT_PKG}/internal/app/build.Version=${VERSION} -X ${PROJECT_PKG}/internal/app/build.CommitHash=${COMMIT_HASH} -X ${PROJECT_PKG}/internal/app/build.BuildDate=${BUILD_DATE}
MOCKS_DESTINATION=test/mocks
.PHONY: mocks

run-orkestrator:
	go run -v ./back-end/orkestrator/cmd/app/main.go serve

run-agent:
	go run -v ./back-end/agent/cmd/app/main.go serve 0

start-docker-compose-test:
	docker-compose -f docker-compose-test.yml up -d

stop-docker-compose-test:
	docker-compose -f docker-compose-test.yml down

test-unit:
	go test -v -cover ./...

test-integration:
	go test -v -tag=integration ./test/...

test-all:
	$(MAKE) test-unit
	$(MAKE) test-integration

.PHONY: build
build-agent:
	go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o agent_app ./back-end/agent/cmd/app

.PHONY: build
build-orkestrator:
	go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o orkestrator_app ./back-end/orkestrator/cmd/app

gen:
	go generate ./...

swagger:
	swag init --parseDependency -g cmd/app/main.go --output=./api

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/grpc/schema/*.proto


install-tools:
	go get -u github.com/onsi/ginkgo/ginkgo
	go install github.com/swaggo/swag/cmd/swag

gen-keys:
	mkdir -p cert
	openssl ecparam -name prime256v1 -genkey -noout -out cert/ec-prime256v1-priv-key.pem
	openssl ec -in cert/ec-prime256v1-priv-key.pem -pubout > cert/ec-prime256v1-pub-key.pem

migrations-up:
	migrate -source file://db/migrations -database postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable up

migrations-down:
	migrate -source file://db/migrations -database postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable down

migrations-up-locally:
	migrate -source file://db/migrations -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable up

migrations-down-locally:
	migrate -source file://db/migrations -database postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable down
