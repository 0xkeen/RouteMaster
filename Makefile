.PHONY: clean build

build:
	go build -o bin/route-master ./cmd

test:
	bash ./scripts/gotest.sh

clean:
	@rm -rf out/
