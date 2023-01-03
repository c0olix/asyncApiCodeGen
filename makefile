install:
	go install github.com/c0olix/asyncApiCodeGen
generate-go:
	asyncApiCodeGen generate -c true -i generator/test-spec/test-spec.yaml -o gensrc/out.gen.go -f mqtt -p events
generate-help:
	asyncApiCodeGen generate -h
help:
	asyncApiCodeGen -h
compile:
	GOOS=windows GOARCH=amd64 go build -o target/asyncApiCodeGen-amd64.exe main.go
	GOOS=darwin GOARCH=amd64 go build -o target/asyncApiCodeGen-amd64-darwin main.go
	GOOS=linux GOARCH=amd64 go build -o target/asyncApiCodeGen-amd64-linux main.go