PROTO=app/users_api/proto

protoc --proto_path=${PROTO} --go_out=${PROTO} --go-grpc_out=${PROTO} \
  users-api.proto

