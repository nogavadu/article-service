LOCAL_BIN:=$(CURDIR)/bin

run:
	docker compose up --build --remove-orphans

stop:
	docker compose down

clear:
	docker compose down -v --rmi local

install-deps:
#	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.24.3

#get-deps:
#	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
#	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

make migration-create:
	$(LOCAL_BIN)/goose -dir "./migrations" create $(NAME) sql