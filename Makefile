TEST_PATH := $(
	if $(word 2,$(MAKECMDGOALS)),internal/$(word 2,$(MAKECMDGOALS))/...,...
)

test:
	go test ./$(TEST_PATH)

repl:
	go run main.go

%:
	@:
