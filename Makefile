.PHONY: build run wire tidy test web web-dev web-build clean

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server -conf configs/config.yaml

# 重新生成依赖注入代码（需: go install github.com/google/wire/cmd/wire@latest）
wire:
	wire ./cmd/server

tidy:
	go mod tidy

test:
	go test ./...

# 前端开发（Vite 热更新，/api 代理到 localhost:8080）
web-dev:
	cd web && npm run dev

# 前端构建（产物 web/dist，由后端静态托管）
web-build:
	cd web && npm run build

web-install:
	cd web && npm install

clean:
	rm -rf bin data web/dist
