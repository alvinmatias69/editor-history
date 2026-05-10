##
# Editor-History
#
# @file
# @version 0.1

.PHONY: all
all:
	go build -o binary cmd/main.go

.PHONY: dev
dev:
	go run cmd/main.go

# end
