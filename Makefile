.PHONY: build clean upload download

build:
	go build -o bin/go-crypt-sync .	

clean:
	rm -rf bin/go-crypt-sync

upload:
	go run main.go encrypt dev_tools/github_recovery_codes.md github_recovery_codes.enc

download:
	go run main.go download encrypted-files-home github_recovery_codes.enc --decrypt true