# Pinned protobuf toolchain. Bump ARGs together when regenerating stubs.
FROM golang:1.27-alpine

ARG PROTOC_VERSION=29.3
ARG PROTOC_GEN_GO_VERSION=1.36.6
ARG PROTOC_GEN_GO_GRPC_VERSION=1.5.1
ARG TARGETARCH

RUN apk add --no-cache curl unzip \
	&& case "${TARGETARCH}" in \
		amd64) PROTOC_ZIP=protoc-${PROTOC_VERSION}-linux-x86_64.zip ;; \
		arm64) PROTOC_ZIP=protoc-${PROTOC_VERSION}-linux-aarch_64.zip ;; \
		*) echo "unsupported TARGETARCH=${TARGETARCH}" >&2; exit 1 ;; \
	esac \
	&& curl -fsSL \
		"https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}" \
		-o /tmp/protoc.zip \
	&& unzip -q /tmp/protoc.zip -d /usr/local \
	&& rm /tmp/protoc.zip \
	&& go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GEN_GO_VERSION} \
	&& go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${PROTOC_GEN_GO_GRPC_VERSION}

WORKDIR /src
COPY docker/proto-gen-entrypoint.sh /usr/local/bin/proto-gen-entrypoint.sh
RUN chmod +x /usr/local/bin/proto-gen-entrypoint.sh

ENTRYPOINT ["proto-gen-entrypoint.sh"]
