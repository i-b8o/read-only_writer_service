APP_BIN = app/build/app

lint:
	golangci-lint run

build: clean $(APP_BIN)

$(APP_BIN):
	go build -o $(APP_BIN) ./app/cmd/main.go

clean:
	rm -rf ./app/build || true

git:
	git add .
	git commit -a -m "$m"
	git push -u origin main
mod:
	cd app;go mod tidy