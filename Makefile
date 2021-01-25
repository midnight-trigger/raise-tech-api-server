gen-api-routing:
	go build -o dev_tools/tools/api_routing_generator dev_tools/api_routing_generator/main.go
	./dev_tools/tools/api_routing_generator ./api/routes.go

gen: gen-api-routing
	go generate ./...

vendor:
	go mod vendor -v

tests:
	go test ./api/definition -v
	go test ./api/domain -v
