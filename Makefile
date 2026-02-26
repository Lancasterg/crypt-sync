.PHONY: build clean

build:
	go build -o bin/go-crypt-sync .	

clean:
	rm -rf bin/go-crypt-sync

# encrypt:
# 	go run 