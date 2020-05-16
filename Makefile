run:
	go run *.go
linux:
	mkdir -p bin
	packr2 clean
	GOOS=linux GOARCH=amd64 packr2 build -ldflags "-s -w" -o "bin/trianglovers_linux" *.go
	packr2 clean
windows:
	mkdir -p bin
	packr2 clean
	GOOS=windows GOARCH=amd64 packr2 build -ldflags "-s -w" -o "bin/trianglovers_windows.exe" *.go
	packr2 clean
darwin:
	mkdir -p bin
	packr2 clean
	rm -rf bin/Trianglovers.app
	rm bin/TriangloversMac.zip
	GOOS=darwin GOARCH=amd64 packr2 build -ldflags "-s -w" -o "bin/trianglovers_mac" *.go
	appify -name "Trianglovers" bin/trianglovers_mac
	mv Trianglovers.app bin
	cd bin && zip TriangloversMac -r ./Trianglovers.app
	packr2 clean
