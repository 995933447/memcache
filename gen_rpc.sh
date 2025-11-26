if [ $# -eq 0 ]; then
  echo "not service input"
  exit 0
fi

mkdir -p $1

protoc --go_out=./$1\
  --go-grpc_out=./$1\
  --easymicro-client_out=./$1\
  --easymicro-server_out=./$1\
  --go_opt=paths=source_relative\
  --go-grpc_opt=paths=source_relative\
  --easymicro-client_opt=paths=source_relative\
  --easymicro-server_opt=paths=source_relative\
  --proto_path=./proto\
  --proto_path=../easymicro/proto\
  $1.proto