PREFIX = /usr/local
BINDIR = $(PREFIX)/bin
TARGET = darakht

.SILENT:
default:
	go run cmd/$(TARGET)/main.go -f misc/sample-35b.txt -c 8 | jq

all: build

run:
	go run cmd/$(TARGET)/main.go

rng:
	g++ misc/src.cpp -o rngF

test:
	go test pkg/merkletree/merkletree_test.go -v

build:
	go build -o $(TARGET) cmd/darakht/main.go

install: all
	install -d $(BINDIR)
	install -m 755 $(TARGET) $(BINDIR)

uninstall:
	rm -f $(BINDIR)/$(TARGET)
