run:
	go run *.go
linux:
	mkdir -p bin
	GOOS=linux GOARCH=amd64 packr2 build -ldflags "-s -w" -o "bin/trianglovers_linux" *.go
	packr clean
windows:
	mkdir -p bin
	GOOS=windows GOARCH=amd64 packr2 build -ldflags "-s -w" -o "bin/trianglovers_windows.exe" *.go
	packr clean
darwin:
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 packr2 build -ldflags "-s -w" -o "bin/trianglovers_mac" *.go
	packr clean
