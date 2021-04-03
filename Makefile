main.wasm:
	GOARCH=wasm GOOS=js go build -o public/wasm/main.wasm go/main.go
