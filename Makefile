NAME=caza
BIN_NAME=caza-bin
DOCKER_REGISTRY=docker.io
DOCKER_REPOSITORY=jrroman
DIR=$(shell cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)
COMPILER_FLAGS="-O2 -g -Wall -Werror"
COMPILER=clang-11
EBPF_DIR=pkg/ebpf

$(NAME): build

generate:
	BPF_CLANG=$(COMPILER) BPF_CFLAGS=$(COMPILER_FLAGS) go generate "$(DIR)/..."

build: generate
	CGO_ENABLED=0 go build -o $(BIN_NAME) $(DIR)/cmd/caza/...

test:
	go test -v ./...

image:
	docker build -t $(DOCKER_REPOSITORY)/$(NAME):latest .

image-push: image
	docker push $(DOCKER_REPOSITORY)/$(NAME):latest

clean:
	rm $(EBPF_DIR)/bpf_bpfel* $(EBPF_DIR)/bpf_bpfeb* $(BIN_NAME)
