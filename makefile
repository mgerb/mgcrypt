linux:
	GOOS=linux GOARCH=amd64 go build -o ./dist/linux/mgcrypt ./main.go

mac:
	GOOS=darwin GOARCH=amd64 go build -o ./dist/mac/mgcrypt ./main.go
	
windows:
	GOOS=windows GOARCH=386 go build -o ./dist/windows/mgcrypt.exe -ldflags="-X main.version=${VERSION}" ./main.go

clean:
	rm -rf ./dist

all: linux mac windows
