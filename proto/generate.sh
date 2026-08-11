#!/usr/bin/env bash
# Build the pinned protobuf toolchain and regenerate or verify every contract
# shipped by the standalone candacelib module. No host Go/protoc installation
# is used.
set -euo pipefail

proto_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
module_root="$(cd "${proto_dir}/.." && pwd)"
mode="${1:-write}"

case "${mode}" in
  write | check) ;;
  *)
    echo "usage: $0 [write|check]" >&2
    exit 2
    ;;
esac

toolchain_image="${CANDACELIB_PROTO_TOOLCHAIN_IMAGE:-candacelib/proto-codegen:go1.26.5-protoc35.1}"

docker build \
  --platform linux/amd64 \
  --file "${proto_dir}/Dockerfile.codegen" \
  --tag "${toolchain_image}" \
  "${proto_dir}"

docker run --rm \
  --platform linux/amd64 \
  --user "$(id -u):$(id -g)" \
  --volume "${module_root}:/workspace" \
  "${toolchain_image}" \
  ./proto/generate-in-container.sh "${mode}"
