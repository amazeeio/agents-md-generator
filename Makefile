.PHONY: all clean serve watch dev

all: dist/main.wasm dist/wasm_exec.js

dist:
	mkdir -p dist

dist/main.wasm: main.go | dist
	GOOS=js GOARCH=wasm go build -o dist/main.wasm

dist/wasm_exec.js: | dist
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" dist/

clean:
	rm -f dist/main.wasm dist/wasm_exec.js

serve: all
	@echo "Serving at: http://localhost:8080"
	cd dist && python3 -m http.server 8080

watch:
	@echo "Watching for changes..."
	@while true; do \
		make all > /dev/null; \
		sleep 1; \
	done

dev: all
	@make watch & \
	WATCH_PID=$$!; \
	trap "kill $$WATCH_PID 2>/dev/null" EXIT INT TERM; \
	python3 dev_server.py
