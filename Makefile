COVERAGE_OUT=coverage.out
install:
	@go install cmd/tagger/tagger.go

test:
	@go test ./...

cov:
	@go test -cover -coverprofile=${COVERAGE_OUT} -coverpkg=./... ./...

cov/view:
	@go tool cover -html=${COVERAGE_OUT}
