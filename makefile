build:
	@mkdir -p build
	@go build -o build/tezos-delegation-service .

test:
	@go clean -testcache && go test ./... -cover

run:
	@go run .