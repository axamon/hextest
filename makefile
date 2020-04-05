APPNAME=hextest
VERSION=0.4.4
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOTOOL=$(GOCMD) tool
GOPPROF=$(GOTOOL) pprof
BINARY_NAME=main
BINARY_UNIX=$(BINARY_NAME)_unix

HOSTPORT=5001

all: test build
build: 
	CGO_ENABLED=0 $(GOBUILD) -ldflags="-w -s -X main.Version=$(VERSION)" -a -installsuffix cgo -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v -cover ./...
testfile: test
	$(GOTEST) -v -cover ./scan -o testscan
profile: test-file
	./testscan --test.v --test.cpuprofile profili/cpu.pprof
	$(GOPPROF) --pdf eseguibili/goscanner-linux profili/cpu.pprof > profili/cpu.pdf
dockerimage: test build
	#CGO_ENABLED=0 $(GOBUILD) -ldflags="-w -s" -a -installsuffix cgo -o main -v
	podman build -t $(APPNAME):$(VERSION) -f Dockerfile
dockerrun:
	podman run --rm -d -it --name $(APPNAME)_$(HOSTPORT) -p $(HOSTPORT):3000 $(APPNAME):$(VERSION) -redis 192.168.1.2:6379
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
