all: constant.go
	go build -o build/suexec bin/suexec.go
test: 
	go test
constant.go:
	ruby source_generate.rb > constant.go
clean:
	rm -f constant.go build/suexec
