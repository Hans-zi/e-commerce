server:
	go run cmd/server/main.go

redis:
	docker start redis7

kafka:
	docker start kafka
.PHONY: server redis kafka