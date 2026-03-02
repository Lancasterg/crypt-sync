.PHONY: build clean upload download,= test

build:
	go build -o bin/go-crypt-sync .	

clean:
	rm -rf bin/go-crypt-sync

upload:
	go run main.go encrypt dev_tools/github_recovery_codes.md github_recovery_codes.enc

download:
	go run main.go download encrypted-files-home github_recovery_codes.enc --decrypt true

list:
	go run main.go list

append:
	go run main.go append SomeService my-username mypassword --recovery=hello:goodbye --recovery=smile:yay --file-name=dev_tools/decrypted.json

test: 
	go test -v ./cmd