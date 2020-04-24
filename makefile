build:
	env CGO_ENABLED=0 GOOS=linux go build -o batch-linux cmd/main.go
	env CGO_ENABLED=0 GOOS=darwin go build -o batch-mac cmd/main.go
	env CGO_ENABLED=0 GOOS=windows go build -o batch.exe cmd/main.go

