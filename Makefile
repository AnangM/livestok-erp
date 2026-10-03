.PHONY: test test-cover test-cover-html clean

test:
	go test -v -race ./...

test-cover:
	go test -v -race -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

test-cover-html:
	go test -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -f coverage.out coverage.html
