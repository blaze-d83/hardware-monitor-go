# Go-related variables
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOFMT = $(GOCMD) fmt
GOTEST = $(GOCMD) test 

# Directories
BINDIR = bin
STATICDIR = static
CMD_DIR = cmd
MAIN_FILE = $(CMD_DIR)/main.go
BINARY_NAME = main

# Targets
.PHONY: all build clean test fmt run

all: build

build: fmt
	$(GOBUILD) -o $(BINDIR)/$(BINARY_NAME) $(MAIN_FILE)

clean:
	$(GOCLEAN)
	rm -f $(BINDIR)/$(BINARY_NAME)

run: build
	$(BINDIR)/$(BINARY_NAME)

test:
	$(GOTEST) ./...

fmt:
	$(GOFMT) ./...

