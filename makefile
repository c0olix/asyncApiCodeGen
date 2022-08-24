install:
	go install github.com/c0olix/asyncApiCodeGen
generate-go:
	asyncApiCodeGen generate generator/test-spec/account-service.yaml gensrc/out.gen.go
compile:
	GOOS=windows GOARCH=amd64 go build -o target/asyncApiCodeGen-amd64.exe main.go
	GOOS=darwin GOARCH=amd64 go build -o target/asyncApiCodeGen-amd64-darwin main.go
	GOOS=linux GOARCH=amd64 go build -o target/asyncApiCodeGen-amd64-linux main.go