.PHONY: docker
docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f ge3m0r/webook:v0.0.1
	@docker build -t ge3m0r/webook:v0.0.1 .
