all: tmp

run:
	go run cmd/test/main.go

rng:
	g++ misc/src.cpp -o rngF

tmp:
	go run cmd/test/main.go misc/sample-35b.txt

build:
	go build -o merkle cmd/test/main.go

clean:
	rm -rf bin/
