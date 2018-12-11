linux:
	GOOS=linux GOARCH=amd64 go build -o ./dist/mgcrypt-linux ./main.go

mac:
	GOOS=darwin GOARCH=amd64 go build -o ./dist/mgcrypt-mac ./main.go
	
windows:
	GOOS=windows GOARCH=386 go build -o ./dist/mgcrypt-windows.exe ./main.go

clean:
	rm -rf ./dist

zip:
	zip -r release.zip dist/*

all: linux mac windows zip
