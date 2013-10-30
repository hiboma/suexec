all:
	go build -o build/suexec bin/suexec.go
test: 
	go test
