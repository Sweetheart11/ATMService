build:
	@go build -o bin/timetracker cmd/main.go 

run: build
	@./bin/timetracker 