USERS_API_PROTO=app/users_api/proto

protoc --proto_path=${USERS_API_PROTO} --go_out=${USERS_API_PROTO} --go-grpc_out=${USERS_API_PROTO} \
  users-api.proto

FIND_NUMBER_POSITION_API_PROTO=app/find_number_position_api/proto

protoc --proto_path=${FIND_NUMBER_POSITION_API_PROTO} --go_out=${FIND_NUMBER_POSITION_API_PROTO} --go-grpc_out=${FIND_NUMBER_POSITION_API_PROTO} \
  find-number-position-api.proto

GATEWAY_PROTO=internal/gateway/proto

cp ${USERS_API_PROTO}/users-api.proto ${GATEWAY_PROTO}/users-api.proto

protoc --proto_path=${GATEWAY_PROTO} --go_out=${GATEWAY_PROTO} --go-grpc_out=${GATEWAY_PROTO} \
  users-api.proto

cp ${FIND_NUMBER_POSITION_API_PROTO}/find-number-position-api.proto ${GATEWAY_PROTO}/find-number-position-api.proto

protoc --proto_path=${GATEWAY_PROTO} --go_out=${GATEWAY_PROTO} --go-grpc_out=${GATEWAY_PROTO} \
  find-number-position-api.proto