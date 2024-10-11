all: tmp

run:
	go run cmd/darakht/main.go

rng:
	g++ misc/src.cpp -o rngF

.SILENT:
tmp:
	go run cmd/darakht/main.go misc/sample-35b.txt

test:
	go test pkg/merkletree/merkletree_test.go -v

build:
	go build -o merkle cmd/test/main.go
