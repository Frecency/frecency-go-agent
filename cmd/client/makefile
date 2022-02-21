
# Go parameters
GOCMD=go
SHRINKFLAGS=-ldflags="-s -w"
DEBUGTAGS=-tags DEBUG
GOBUILD=$(GOCMD) build $(SHRINKFLAGS)
GOBUILDDEBUG=$(GOCMD) build $(DEBUGTAGS)
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=client
BINARY_LINUX=$(BINARY_NAME)_linux
BINARY_WIN=$(BINARY_NAME)_win.exe
BINARY_MAC=$(BINARY_NAME)_mac
DEBUG_BINARY_LINUX=debug_$(BINARY_NAME)_linux
DEBUG_BINARY_WIN=debug_$(BINARY_NAME)_win.exe
DEBUG_BINARY_MAC=debug_$(BINARY_NAME)_mac

prod:
	make build-linux
	make build-windows
	make build-mac
debug:
	make build-linux-debug
	make build-windows-debug
	make build-mac-debug
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_LINUX)
	rm -f $(BINARY_WIN)
	rm -f $(BINARY_MAC)
	rm -f $(DEBUG_BINARY_LINUX)
	rm -f $(DEBUG_BINARY_WIN)
	rm -f $(DEBUG_BINARY_MAC)

# build all platforms
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -v