
.PHONY install:
install:
	go install .

.PHONY test:
test:
	go test ./...

.PHONY generate:
generate:
	go generate ./...


