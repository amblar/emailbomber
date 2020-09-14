GO := go
VERSION := v0.0.1
INSTALL_DIR := /usr/bin

ifeq ($(shell uname), Darwin)
	INSTALL_DIR = /usr/local/bin
endif

.PHONY: build release_windows release_linux release_darwin clean test deps

build:
	$(GO) build -v -ldflags "-s -w" && upx Emailbomber

install:
	@if [ ! -d $(INSTALL_DIR) ]; then \
		echo "Unable to locate $(INSTALL_DIR)"; \
		exit 1; \
	fi

	@if [ ! -w $(INSTALL_DIR) ]; then \
		echo "Insufficient permissions, please elevate"; \
		exit 1; \
	fi

	@if [ ! -f emailbomber ]; then \
		$(MAKE) build; \
	fi

	@cp emailbomber $(INSTALL_DIR)/emailbomber && echo "Installed successfully into $(INSTALL_DIR), use emailbomber command to start"; \

release_windows:
	mkdir emailbomber-$(VERSION)-windows-amd64
	env GOOS=windows GOARCH=amd64 $(GO) build -v -o emailbomber-$(VERSION)-windows-amd64/
	zip -9 emailbomber-$(VERSION)-windows-amd64.zip emailbomber-$(VERSION)-windows-amd64/emailbomber.exe
	echo `shasum -a 256 emailbomber-$(VERSION)-windows-amd64.zip` > emailbomber-$(VERSION)-windows-amd64.sha256
	rm -rf emailbomber-$(VERSION)-windows-amd64

release_linux:
	mkdir emailbomber-$(VERSION)-linux-amd64
	env GOOS=linux GOARCH=amd64 $(GO) build -v -o emailbomber-$(VERSION)-linux-amd64/
	tar -cvzf emailbomber-$(VERSION)-linux-amd64.tar.gz emailbomber-$(VERSION)-linux-amd64/emailbomber
	echo `shasum -a 256 emailbomber-$(VERSION)-linux-amd64.tar.gz` > emailbomber-$(VERSION)-linux-amd64.sha256
	rm -rf emailbomber-$(VERSION)-linux-amd64

release_darwin:
	mkdir emailbomber-$(VERSION)-darwin-amd64
	env GOOS=darwin GOARCH=amd64 $(GO) build -v -o emailbomber-$(VERSION)-darwin-amd64/
	tar -cvzf emailbomber-$(VERSION)-darwin-amd64.tar.gz emailbomber-$(VERSION)-darwin-amd64/emailbomber
	echo `shasum -a 256 emailbomber-$(VERSION)-darwin-amd64.tar.gz` > emailbomber-$(VERSION)-darwin-amd64.sha256
	rm -rf emailbomber-$(VERSION)-darwin-amd64

clean:
	$(GO) clean
	rm -rf emailbomber-$(VERSION)
	rm -rf emailbomber-$(VERSION)-windows-amd64
	rm -rf emailbomber-$(VERSION)-linux-amd64
	rm -rf emailbomber-$(VERSION)-darwin-amd64
	rm -f emailbomber-$(VERSION)-windows-amd64.zip
	rm -f emailbomber-$(VERSION)-linux-amd64.tar.gz
	rm -f emailbomber-$(VERSION)-darwin-amd64.tar.gz
	rm -f emailbomber-$(VERSION)-windows-amd64.sha256
	rm -f emailbomber-$(VERSION)-linux-amd64.sha256
	rm -f emailbomber-$(VERSION)-darwin-amd64.sha256

test:
	$(GO) test ./...

deps:
	$(GO) get -v -u github.com/fatih/color