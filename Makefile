DIR=$(shell cd "$(dirname "${BASH_SOURCE[0]}")" &> /dev/null && pwd)
COMPILER_FLAGS="-O2 -g -Wall -Werror"
COMPILER=clang-12

build:
	BPF_CLANG=$(COMPILER) BPF_CFLAGS=$(COMPILER_FLAGS) go generate "$(DIR)/..."
