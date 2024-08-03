run:
	clear
	go run ./cmd/api -cors-trusted-origins="http://localhost:9000 http://localhost:9001"