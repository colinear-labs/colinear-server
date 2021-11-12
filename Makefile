SHELL = /bin/sh
.DEFAULT_GOAL := release
COMMIT_HASH := $(shell git rev-parse --short HEAD)
WIDGET_DIR := ./x-payment-widget

build: 
	@mkdir -p bin && \
	echo Building x-server for all architectures.
	GOOS=linux GOARCH=arm go build -o bin/xserver-linux-arm; \
	GOOS=linux GOARCH=arm64 go build -o bin/xserver-linux-arm64; \
	GOOS=linux GOARCH=386 go build -o bin/xserver-linux-386; \
	GOOS=linux GOARCH=amd64 go build -o bin/xserver-linux-amd64; \
	GOOS=darwin GOARCH=amd64 go build -o bin/xserver-darwin-amd64; \
	GOOS=darwin GOARCH=arm64 go build -o bin/xserver-darwin-arm64; \

build-widget:
	@cd x-payment-widget && \
	echo Building payment widget. && \
	yarn && \
	yarn build

clean:
	@rm -rf bin
	@rm -rf release
	@rm -rf widget

dev: build-widget
	@cp -r ${WIDGET_DIR}/public ./widget

release: build-widget build
	@mkdir -p release
	@mkdir -p release/x-server-${COMMIT_HASH}
	@mv bin/* release/x-server-${COMMIT_HASH}
	@rmdir bin
	@cp -r ${WIDGET_DIR}/public release/x-server-${COMMIT_HASH}/widget

release-docker: release
