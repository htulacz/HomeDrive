clean:
	rm -fr upload HomeDrive*

all:
	go build
	cd site
	npm run build
	cd ..
windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o HomeDrive.exe server.go