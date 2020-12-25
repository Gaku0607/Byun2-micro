.PHONY: run-member
run-member:
	go run ./member-api/main.go

.PHONY: gen-member
gen-member:
	protoc --proto_path=member-api/proto model.proto --go_out=plugins=grpc:member-api/proto/pb --micro_out=member-api/proto/pb
	protoc -I/usr/local/include --proto_path=member-api/proto member.proto --go_out=plugins=grpc:member-api/proto/pb --micro_out=member-api/proto/pb

.PHONY: clean-member
clean-member:
	rm ./member-api/proto/pb/*.go
	