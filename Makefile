HOSTNAME=registry.terraform.io
NAMESPACE=bluechi
NAME=bluechi
BINARY=terraform-provider-${NAME}
VERSION=0.1.0
OS_ARCH=linux_amd64
INSTALLDIR=~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

SOURCES = $(shell find . -path './.*' -prune -o \( \( -name '*.go' \) -a ! -name '*_test.go' \) -print)

.PHONY: default
default: all

.PHONY: all
all: install

.PHONY: build
build: ${BINARY}

.PHONY: ${BINARY}
${BINARY}: ${SOURCES} go.mod go.sum
	go build -o ${BINARY}

.PHONY: install
install: ${BINARY}
	mkdir -p ${INSTALLDIR}
	cp ${BINARY} ${INSTALLDIR}

uninstall:
	rm -f ${INSTALLDIR}/${BINARY}

clean: uninstall
	rm -f ${BINARY}

test:
	bash container/container-setup.sh start bluechi
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
	bash container/container-setup.sh stop
