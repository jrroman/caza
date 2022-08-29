#!/usr/bin/env bash

DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

# Version of libbpf to fetch headers from
LIBBPF_VERSION=0.6.1

# The headers we want
prefix=libbpf-"$LIBBPF_VERSION"
headers=(
    "$prefix"/LICENSE.BSD-2-Clause
    "$prefix"/src/bpf_endian.h
    "$prefix"/src/bpf_helper_defs.h
    "$prefix"/src/bpf_helpers.h
    "$prefix"/src/bpf_tracing.h
)

# Fetch libbpf release and extract the desired headers
curl -sL "https://github.com/libbpf/libbpf/archive/refs/tags/v${LIBBPF_VERSION}.tar.gz" | \
    tar -xz --xform='s#.*/##' "${headers[@]}"

# create vmlinux header file based on host kernel
sudo bpftool btf dump file /sys/kernel/btf/vmlinux format c > "$DIR/vmlinux.h"
