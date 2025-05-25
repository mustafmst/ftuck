all: clean build install-user

build:
	@echo "BUILD"
	@mkdir -p dist 
	@go build -o dist/

clean:
	@echo "CLEAN"
	@rm -rf ./dist

install-user:
	@echo "USER WIDE INSTALL TO ~/.local/bin"
	@mkdir -p ~/.local/bin
	@cp -r ./dist/ftuck ~/.local/bin/ 
	@echo "command: ftuck"
	@echo "To enable usage make sure that ~/.local/bin is in Your PATH"

.PHONY: test
test:
	go test ./...
