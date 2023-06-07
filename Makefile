
all: test

.PHONY:test
test:
	go test -v ./test/...

clean:
	rm -rf ./test/logs