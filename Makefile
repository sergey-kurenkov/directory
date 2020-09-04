.PHONY: all test testfv build run-linter demo

all: test build

test:
	go test -count=1 ./internal

build:
	go build ./cmd/query_dir

testfv:
	go test -failfast -v -count=1 ./internal

run-linter:
	golangci-lint run ./internal

demo: build
	@echo ""

	@echo "Finding a common manager for Ann and Julia:"
	@./query_dir -file=./test/test_org2.json -first=Ann -second=Julia
	@echo ""

	@echo "Finding a common manager for Ann and Bob:"
	@./query_dir -file=./test/test_org2.json -first=Bob -second=Ann	
	@echo ""

	@echo "Finding a common manager for Joseph and Kate:"
	@./query_dir -file=./test/test_org2.json -first=Joseph -second=Kate
	@echo ""

	@echo "Finding a common manager for Joseph and Donald:"
	@./query_dir -file=./test/test_org2.json -first=Joseph -second=Donald
	@echo ""

	@echo "Finding a common manager for Joseph and Ann:"
	@./query_dir -file=./test/test_org2.json -first=Joseph -second=Ann
	@echo ""

	@echo "Finding a common manager for Monica and Monica:"
	@./query_dir -file=./test/test_org2.json -first=Monica -second=Monica
	@echo ""

	@echo "Finding a common manager for Jane and Jane:"
	@./query_dir -file=./test/test_org2.json -first=Jane -second=Jane
	@echo ""

	@echo "Finding a common manager for Bill and Bill:"
	@./query_dir -file=./test/test_org2.json -first=Bill -second=Bill
	@echo ""

	@echo "Finding a common manager for John and John:"
	@./query_dir -file=./test/test_org2.json -first=John -second=John
	@echo ""

	@echo "Finding a common manager for Charlie and Mia:"
	@./query_dir -file=./test/test_org3.json -first=Charlie -second=Mia
	@echo ""

	@echo "Finding a common manager for Charlie and Boris:"
	@./query_dir -file=./test/test_org3.json -first=Charlie -second=Boris
	@echo ""

	@echo "Finding a common manager for Mike and Kirk:"
	@./query_dir -file=./test/test_org3.json -first=Mike -second=Kirk
	@echo ""
