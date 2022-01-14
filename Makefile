GOBIN ?= ${GOPATH}/bin

all: addon

fmt:
	go fmt ./...

vet:
	go vet ./...

addon: fmt vet
	go build -ldflags="-w -s" -o bin/addon github.com/konveyor/tackle-addons-analysis-windup/cmd
