build:
	go build -o ./bin/server/server ./cmd/server

api_go:
	protoc -I ./api/proto/v1/ \
	-I ${GOPATH}/src \
	-I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway \
	--go_out=plugins=grpc:pkg/api/v1 \
	economy_service.proto

api_gateway:
	protoc -I ./api/proto/v1/ \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway \
	--grpc-gateway_out=logtostderr=true:pkg/api/v1 \
	economy_service.proto

api_swagger:
	protoc -I ./api/proto/v1/ \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway \
	--swagger_out=logtostderr=true:api/swagger/v1 \
	economy_service.proto

api: api_go api_gateway api_swagger

mocks:
	cd pkg && mockery -all
