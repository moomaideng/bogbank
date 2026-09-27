#!/bin/sh
set -eu

MODULE="${MODULE:-github.com/moomaideng/bogbank}"

# Collect under /src so Windows hosts never need find on PATH.
set --
for f in $(find services -name '*.proto' | sort); do
	set -- "$@" "$f"
done

if [ "$#" -eq 0 ]; then
	echo "no .proto files under services/" >&2
	exit 1
fi

protoc \
	--go_out=. --go_opt=module="$MODULE" \
	--go-grpc_out=. --go-grpc_opt=module="$MODULE" \
	"$@"
