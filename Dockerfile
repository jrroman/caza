FROM golang:bullseye

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /workspace

COPY . .

RUN apt update && \
    apt install -y \
        clang-11 \
        libclang-11-dev

CMD ./start.sh
