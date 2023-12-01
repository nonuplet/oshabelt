buf-lint:
	cd api && buf lint
buf-gen:
	cd api && buf generate

dev:
	cd backend && go run cmd/main.go