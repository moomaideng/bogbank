MODULE := github.com/moomaideng/bogbank
PROTOS := $(shell find services -name '*.proto')

# protoc-gen-go and protoc-gen-go-grpc must be on PATH.
.PHONY: proto
proto:
	protoc \
		--go_out=. --go_opt=module=$(MODULE) \
		--go-grpc_out=. --go-grpc_opt=module=$(MODULE) \
		$(PROTOS)

.PHONY: dev-ledger
dev-ledger:
	docker compose up -d
	go run ./services/ledger serve --config services/ledger/config.yaml

.PHONY: down
down:
	docker compose down
