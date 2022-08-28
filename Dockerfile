# use golang:bullseye for building our program
FROM golang:bullseye as builder

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /workspace

COPY . .

# Install dependencies, clean apt and build program
RUN apt update && \
    apt install -y \
        clang-11 \
        libclang-11-dev \
        && apt autoremove \
        && rm -rf /var/lib/apt/lists/* \
        && make build

# The image we will be running for the container
FROM amazonlinux:2

WORKDIR /workspace

COPY --from=builder /workspace/caza ./

CMD ./caza
