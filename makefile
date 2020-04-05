APPNAME=hextest
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

all: test build
build: 
	CGO_ENABLED=0 $(GOBUILD) -ldflags="-w -s" -a -installsuffix cgo -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v -cover ./...
testfile: test
	$(GOTEST) -v -cover ./scan -o testscan
profile: test-file
	./testscan --test.v --test.cpuprofile profili/cpu.pprof
	$(GOPPROF) --pdf eseguibili/goscanner-linux profili/cpu.pprof > profili/cpu.pdf
dockerimage: test
	CGO_ENABLED=0 $(GOBUILD) -ldflags="-w -s" -a -installsuffix cgo -o main
	podman build -t $(APPNAME):latest -f Dockerfile
dockerrun:
	podman run --rm -d -it --name $(APPNAME) -p 5001:3000 $(APPNAME):latest -redis 192.168.1.2:6379
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
