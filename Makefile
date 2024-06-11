
# go params
GOCMD=go

# normal entry points
	
update:
	clear 
	@echo "updating dependencies..."
	@go get -u -t ./...
	@go mod tidy 

build:
	clear 
	@echo "building..."
	@$(GOCMD) build .
	
test:
	clear
	@echo "testing Jobnimbus..."
	@$(GOCMD) test -run TestFirst ./...

test-second:
	clear
	@echo "test Jobnimbus second level functions..."
	@$(GOCMD) test -run TestSecond ./...

test-third:
	clear
	@echo "test Jobnimbus third level functions..."
	@$(GOCMD) test -run TestThird ./...
