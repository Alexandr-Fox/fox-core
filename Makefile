all: build

run: generate
	go run main.go

build: generate
	go build -o build/

install:
	cp build/* /usr/bin/

generate:
	go generate

clean:
	rm -r build