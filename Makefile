clean:
	rm -fr upload HomeDrive*

all:
	go build

windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o HomeDrive.exe server.go