build:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 packr build -ldflags "-s -w" -o "bin/trianglovers_linux" *.go
windows:
	GOOS=windows GOARCH=amd64 packr build -ldflags "-s -w" -o "bin/trianglovers_windows.exe" *.go
darwin:
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 packr build -ldflags "-s -w" -o "bin/trianglovers_mac" *.go; \
