all:
	go build -o build/suexec bin/suexec.go
test: 
	go test
	cd env    && go test
	cd passwd && go test