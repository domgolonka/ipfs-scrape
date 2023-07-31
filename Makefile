docker-build:
	docker build -t blockparty:latest .
	docker tag blockparty:latest

run-api:
	go run cmd/api/main.go

run-scraper:
	go run cmd/scraper/main.go

build-%-static: cmd/%
	-rm -r bin
	mkdir -p bin
	GOOS=linux CGO_ENABLED=0 go build  -buildvcs=false -o  ./bin/$(*F) ./cmd/$(*F)