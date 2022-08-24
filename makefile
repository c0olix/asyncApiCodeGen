install:
	go install github.com/c0olix/asyncApiCodeGen
generate-go:
	asyncApiCodeGen generate generator/test-spec/account-service.yaml gensrc/out.gen.go