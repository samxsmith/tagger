COVERAGE_OUT=coverage.out
test:
	@go test ./...

cov:
	@go test -cover -coverprofile=${COVERAGE_OUT} -coverpkg=./... ./...

cov/view:
	@go tool cover -html=${COVERAGE_OUT}
