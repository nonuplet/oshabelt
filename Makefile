buf-lint:
	cd api && buf lint
buf-gen:
	cd api && buf generate
dev-back:
	cd backend && go run cmd/main.go
dev-front:
	cd frontend && yarn dev