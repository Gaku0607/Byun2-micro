.PHONY: run-member
run-member: 
	go run ./member-srv/main.go

.PHONY: gen-member
gen-member:
	protoc --proto_path=member-srv/proto model.proto --go_out=plugins=grpc:member-srv/pb --micro_out=member-srv/pb
	protoc -I/usr/local/include --proto_path=member-srv/proto member.proto --go_out=plugins=grpc:member-srv/pb --micro_out=member-srv/pb

.PHONY:clean-member
clean-member:
	rm  member-srv/pb/*.go