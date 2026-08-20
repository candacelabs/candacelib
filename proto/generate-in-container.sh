#!/usr/bin/env bash
# Runs inside the pinned code-generation image. Use generate.sh from the host.
set -euo pipefail

proto_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
module_root="$(cd "${proto_dir}/.." && pwd)"
mode="${1:-write}"
case "${mode}" in
  write)
    output_root="${module_root}"
    ;;
  check)
    output_root="$(mktemp -d /tmp/candacelib-proto-output.XXXXXX)"
    ;;
  *)
    echo "usage: $0 [write|check]" >&2
    exit 2
    ;;
esac

test "$(protoc --version)" = "libprotoc 35.1"
test "$(protoc-gen-go --version)" = "protoc-gen-go v1.36.11"

export GOCACHE=/tmp/candacelib-proto-gocache
export GOMODCACHE=/tmp/candacelib-proto-modcache
export GOPATH=/tmp/candacelib-proto-gopath

plugin_dir="$(mktemp -d /tmp/candacelib-liquid-plugin.XXXXXX)"
(cd "${module_root}" && \
  go build -mod=readonly -o "${plugin_dir}/protoc-gen-liquidproto" ./cmd/protoc-gen-liquidproto)

protoc \
  -I "${module_root}" \
  -I /usr/local/include \
  "--go_out=module=github.com/candacelabs/candacelib:${output_root}" \
  liquidproto/v1/refinement.proto

protoc \
  -I "${module_root}" \
  -I /usr/local/include \
  "--go_out=module=github.com/candacelabs/candacelib:${output_root}" \
  "--plugin=protoc-gen-liquidproto=${plugin_dir}/protoc-gen-liquidproto" \
  "--liquidproto_out=module=github.com/candacelabs/candacelib:${output_root}" \
  boundedbuffer/v1/buffer.proto \
  cron/v1/cron.proto

if [[ "${mode}" == check ]]; then
  generated_files=(
    liquidproto/v1/refinement.pb.go
    boundedbuffer/v1/buffer.pb.go
    boundedbuffer/v1/buffer_liquid.pb.go
    cron/v1/cron.pb.go
    cron/v1/cron_liquid.pb.go
  )
  for generated_file in "${generated_files[@]}"; do
    diff -u "${module_root}/${generated_file}" "${output_root}/${generated_file}"
  done
fi
