.PHONY: test coverhtml

test:
	go test -v --race .

coverhtml:
	@mkdir -p coverage
	@go test -coverprofile=coverage/cover.out .
	@go tool cover -html=coverage/cover.out -o coverage/coverage.html
	@go tool cover -func=coverage/cover.out | tail -n 1
