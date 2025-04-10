build:
	go build -o ./bin ./cmd

run: build
	./bin/cmd

test: 
	go test ./... 

test_desc: 
	go test -v ./... 