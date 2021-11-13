SHELL = /bin/sh
.DEFAULT_GOAL := release
COMMIT_HASH := $(shell git rev-parse --short HEAD)
WIDGET_DIR := ./x-payment-widget
WEBUI_DIR := ./x-webui

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

build-webui:
	@cd x-webui && \
	echo Building merchant web UI. && \
	yarn && \
	yarn build

clean:
	@rm -rf bin
	@rm -rf release
	@rm -rf widget
	@rm -rf webui

dev: build-widget build-webui
	@mkdir -p ./widget
	@mkdir -p ./webui
	@cp -r ${WIDGET_DIR}/public/* ./widget
	@cp -r ${WEBUI_DIR}/public/* ./webui

release: build-widget build
	@mkdir -p release
	@mkdir -p release/x-server-${COMMIT_HASH}
	@mv bin/* release/x-server-${COMMIT_HASH}
	@rmdir bin
	@mkdir -p release/x-server-${COMMIT_HASH}/widget
	@mkdir -p release/x-server-${COMMIT_HASH}/webui
	@cp -r ${WIDGET_DIR}/public/* release/x-server-${COMMIT_HASH}/widget
	@cp -r ${WEBUI_DIR}/public/* release/x-server-${COMMIT_HASH}/webui

release-docker: release
