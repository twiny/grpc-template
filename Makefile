# generate pkg file
generate:
	buf generate proto/

# remove pkg file
clean:
	rm -rf ./pkg

# test pkg
test:
	go test -cover -race ./...