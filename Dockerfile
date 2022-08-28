FROM golang:bullseye as builder

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /workspace

COPY . .

RUN apt update && \
    apt install -y \
        clang-11 \
        libclang-11-dev \
        && rm -rf /var/lib/apt/lists/*

RUN make build

FROM amazonlinux:2

WORKDIR /workspace

COPY --from=builder /workspace/caza ./

CMD ./caza
