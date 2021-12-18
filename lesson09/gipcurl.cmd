#list services from running server
grpcurl -plaintext -v localhost:9000 list
#list serv declared in myservice.proto file
grpcurl -plaintext -v localhost:9000 list MessageService
#info
grpcurl -plaintext -v localhost:9000 describe MessageService
#info
grpcurl -plaintext -v localhost:9000 describe MessageService.SendMessage

#try service - send json format message using MessageService/SendMessage:
grpcurl -plaintext -d '{"id": "121", "body": "ops message"}' localhost:9000 MessageService/SendMessage