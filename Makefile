build:
	go build -o ./bin ./cmd

run: build
	./bin/cmd

test: 
	go test ./... 

test_desc: 
	go test -v ./... 

up: build
	./bin/cmd localhost:8080 
up1: build
	./bin/cmd localhost:8080 localhost:8081 localhost:8082

up2: build
	./bin/cmd localhost:8081 localhost:8080 localhost:8082

up3: build
	./bin/cmd localhost:8082 localhost:8080 localhost:8081
