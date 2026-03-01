.PHONY: build clean upload download


build:
	go build -o bin/go-crypt-sync .	

clean:
	rm -rf bin/go-crypt-sync

upload:
	go run main.go encrypt dev_tools/test.json hello098.enc

download:
	go run main.go download encrypted-files-home hello123.enc --output dev_tools/test123.json