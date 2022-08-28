NAME=caza

DIR=$(shell cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)
COMPILER_FLAGS="-O2 -g -Wall -Werror"
COMPILER=clang-11

$(NAME): build

generate:
	BPF_CLANG=$(COMPILER) BPF_CFLAGS=$(COMPILER_FLAGS) go generate "$(DIR)/..."

build: generate
	CGO_ENABLED=0 go build -o $(NAME) ./

clean:
	rm bpf_bpfel* bpf_bpfeb* $(NAME)
