wasm:
	GOARCH=wasm GOOS=js go build -ldflags="-s -w" -trimpath -o _web/main.wasm ./cmd/game

itchio-wasm: wasm
	cd _web && \
		mkdir -p ../bin && \
		rm -f ../bin/nebuleet_troopers.zip && \
		zip ../bin/nebuleet_troopers.zip -r main.wasm index.html wasm_exec.js
