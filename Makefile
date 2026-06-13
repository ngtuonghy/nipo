.PHONY: build build-all clean

build:
	cd cli && go build -o ../dist/nipo ./cmd/nipo

build-all:
	@mkdir -p dist
	cd cli && GOOS=windows GOARCH=amd64 go build -o ../dist/nipo-windows-amd64.exe ./cmd/nipo
	cd cli && GOOS=linux GOARCH=amd64 go build -o ../dist/nipo-linux-amd64 ./cmd/nipo
	cd cli && GOOS=linux GOARCH=arm64 go build -o ../dist/nipo-linux-arm64 ./cmd/nipo
	cd cli && GOOS=darwin GOARCH=amd64 go build -o ../dist/nipo-darwin-amd64 ./cmd/nipo
	cd cli && GOOS=darwin GOARCH=arm64 go build -o ../dist/nipo-darwin-arm64 ./cmd/nipo

clean:
	rm -rf dist/
	rm -f cli/nipo.exe cli/nipo
